package sync

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/juju/errors"
	gitignore "github.com/sabhiram/go-gitignore"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// CopyToContainer copies a local folder to a container path
func CopyToContainer(Kubectl kubernetes.Interface, Pod *k8sv1.Pod, Container *k8sv1.Container, LocalPath, ContainerPath string, ExcludePaths []string) error {
	return copyToContainerTestable(Kubectl, Pod, Container, LocalPath, ContainerPath, ExcludePaths, false)
}

func copyToContainerTestable(Kubectl kubernetes.Interface, Pod *k8sv1.Pod, Container *k8sv1.Container, LocalPath, ContainerPath string, ExcludePaths []string, testing bool) error {
	stat, err := os.Lstat(LocalPath)

	if err != nil {
		return errors.Trace(err)
	}

	if stat.IsDir() == false {
		LocalFile := LocalPath
		LocalPath = filepath.Dir(LocalPath)

		files, readDirErr := ioutil.ReadDir(LocalPath)
		if readDirErr != nil {
			return errors.Trace(readDirErr)
		}

		for _, file := range files {
			if filepath.Join(LocalPath, file.Name()) != LocalFile {
				ExcludePaths = append(ExcludePaths, "/"+file.Name())
			}
		}
	}

	s := &SyncConfig{
		Kubectl:      Kubectl,
		Pod:          Pod,
		Container:    Container,
		WatchPath:    getRelativeFromFullPath(LocalPath, ""),
		DestPath:     ContainerPath,
		ExcludePaths: ExcludePaths,
		silent:       true,
		testing:      testing,
	}

	// syncLog = log.GetInstance()
	s.setup()

	err = s.upstream.start()

	if err != nil {
		return errors.Trace(err)
	}

	err = s.upstream.applyCreates([]*fileInformation{
		{
			Name:        "",
			IsDirectory: true,
		},
	})

	if err != nil {
		return errors.Trace(err)
	}

	s.Stop(nil)

	// syncLog = nil

	return nil
}

// We need this function because tar rounds the mtime on the server as well
func roundMtime(mtime time.Time) int64 {
	return mtime.Round(time.Second).Unix()
}

func getRelativeFromFullPath(fullpath string, prefix string) string {
	return strings.TrimPrefix(strings.Replace(strings.Replace(fullpath[len(prefix):], "\\", "/", -1), "//", "/", -1), ".")
}

func pipeStream(w io.Writer, r io.Reader) error {
	buf := make([]byte, 1024, 1024)

	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]

			_, err := w.Write(d)
			if err != nil {
				return errors.Trace(err)
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return errors.Trace(err)
		}
	}
}

func readTill(keyword string, reader io.Reader) (string, error) {
	var output bytes.Buffer
	buf := make([]byte, 0, 512)
	overlap := ""

	for keywordFound := false; keywordFound == false; {
		n, err := reader.Read(buf[:cap(buf)])

		buf = buf[:n]

		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}

			return "", errors.Trace(err)
		}

		// process buf
		if err != nil && err != io.EOF {
			return "", errors.Trace(err)
		}

		lines := strings.Split(string(buf), "\n")

		for index, element := range lines {
			line := ""

			if index == 0 {
				if len(lines) > 1 {
					line = overlap + element
				} else {
					overlap += element
				}
			} else if index == len(lines)-1 {
				overlap = element
			} else {
				line = element
			}

			if line == keyword {
				output.WriteString(line)
				keywordFound = true
				break
			} else if overlap == keyword {
				output.WriteString(overlap)
				keywordFound = true
				break
			} else if line != "" {
				output.WriteString(line + "\n")
			}
		}
	}

	return output.String(), nil
}

func waitTill(keyword string, reader io.Reader) error {
	buf := make([]byte, 0, 512)
	overlap := ""

	for {
		n, err := reader.Read(buf[:cap(buf)])

		buf = buf[:n]

		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}

			return errors.Trace(err)
		}

		// process buf
		if err != nil && err != io.EOF {
			return errors.Trace(err)
		}

		lines := strings.Split(string(buf), "\n")

		for index, element := range lines {
			line := ""

			if index == 0 {
				if len(lines) > 1 {
					line = overlap + element
				} else {
					overlap += element
				}
			} else if index == len(lines)-1 {
				overlap = element
			} else {
				line = element
			}

			if line == keyword || overlap == keyword {
				return nil
			}
		}
	}

	return nil
}

// clean prevents path traversals by stripping them out.
// This is adapted from https://golang.org/src/net/http/fs.go#L74
func clean(fileName string) string {
	return path.Clean(string(os.PathSeparator) + fileName)
}

// dirExists checks if a path exists and is a directory.
func dirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.Trace(err)
}

func deleteSafeRecursive(basepath, relativePath string, fileMap map[string]*fileInformation, removeFiles map[string]*fileInformation, config *SyncConfig) {
	absolutePath := path.Join(basepath, relativePath)
	relativePath = getRelativeFromFullPath(absolutePath, basepath)

	// We don't delete the folder or the contents if we haven't tracked it
	if fileMap[relativePath] == nil || removeFiles[relativePath] == nil {
		config.Logf("[Downstream] Skip delete directory %s\n", relativePath)
		return
	}

	// Delete directory from fileMap
	defer delete(fileMap, relativePath)
	files, err := ioutil.ReadDir(absolutePath)
	if err != nil {
		return
	}

	for _, f := range files {
		filepath := path.Join(relativePath, f.Name())
		absFilepath := path.Join(basepath, filepath)

		if shouldRemoveLocal(absFilepath, fileMap[filepath], config) {
			if f.IsDir() {
				deleteSafeRecursive(basepath, filepath, fileMap, removeFiles, config)
			} else {
				err = os.Remove(absFilepath)
				if err != nil {
					config.Logf("[Downstream] Skip file delete %s: %v", relativePath, err)
				}
			}
		} else {
			config.Logf("[Downstream] Skip delete %s", relativePath)
		}

		delete(fileMap, filepath)
	}

	// This will not remove the directory if there is still a file or directory in it
	err = os.Remove(absolutePath)
	if err != nil {
		config.Logf("[Downstream] Skip delete directory %s, because %s\n", relativePath, err.Error())
	}
}

func compilePaths(excludePaths []string) (gitignore.IgnoreParser, error) {
	if len(excludePaths) > 0 {
		ignoreParser, err := gitignore.CompileIgnoreLines(excludePaths...)

		if err != nil {
			return nil, errors.Trace(err)
		}

		return ignoreParser, nil
	}

	return nil, nil
}

func cleanupSyncLogs() error {
	syncLogName := log.Logdir + "sync.log"
	_, err := os.Stat(syncLogName)

	if err != nil {
		return nil
	}

	// We read the log file and append it to the old log
	data, err := ioutil.ReadFile(syncLogName)

	if err != nil {
		return err
	}

	// Append to syncLog.log.old
	f, err := os.OpenFile(syncLogName+".old", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(data); err != nil {
		return err
	}

	err = os.Remove(syncLogName)

	if err != nil {
		return err
	}

	return nil
}
