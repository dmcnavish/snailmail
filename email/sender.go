package email

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const groupFileSize int64 = 1000000

// ReadFileAndZip given a file name, it reads the file and creates multiple zip files of size groupFileSize
func ReadFileAndZip(fileName string) error {
    if fileName == "" {
        log.Fatal("File name is required to send file!")
    }
    
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer closeFile(f)

	buf := make([]byte, groupFileSize)
	for i := 0; ; i++ {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal("Error reading file: ", err)
		}

		if n == 0 {
			break
		}

		outputFileName := fmt.Sprintf("%s.%d", fileName, i)
		writeFile(outputFileName, buf[:n])
		defer deleteFile(outputFileName)
		err = createZipFile(outputFileName, outputFileName+".zip")
		if err != nil {
			log.Println("Error creating zip", err)
			return err
		}
	}

	return nil
}

func deleteFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Fatal("failed to delete file: ", err)
	}
}

func writeFile(targetFileName string, data []byte) error {
	fileOutput, err := os.OpenFile(targetFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer closeFile(fileOutput)

	if _, err := fileOutput.Write(data); err != nil {
		return err
	}

	return nil
}

func createZipFile(sourceFileName, targetFileName string) error {
	zipfile, err := os.Create(targetFileName)
	if err != nil {
		log.Fatal("error creating zip file: "+sourceFileName, err)
		return err
	}
	defer closeFile(zipfile)

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(sourceFileName)
	if err != nil {
		log.Fatal(err)
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		log.Fatal(err)
		return err
	}

	header.Method = zip.Deflate

	writer, err := archive.CreateHeader(header)
	if err != nil {
		log.Fatal(err)
		return err
	}

	file, err := os.Open(sourceFileName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer closeFile(file)

	_, err = io.Copy(writer, file)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// UnzipAndJoin takes the base name of a collection of zip files, joins, and then extracts them
func UnzipAndJoin(sourceBaseName string) error {
	zipFiles := []string{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, sourceBaseName) && strings.Contains(path, ".zip") && !info.IsDir() {
			zipFiles = append(zipFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Println("zipFiles: ", zipFiles)
  if len(zipFiles) == 0 {
    log.Println("No zip files found! Exiting")
    return nil
  }

	for _, zipFile := range zipFiles {
		err := unzip(zipFile, sourceBaseName)
		if err != nil {
			log.Fatal("error unzipping file: "+ zipFile, err)
			panic(err)
		}
	}

	err = joinFiles(sourceBaseName)
	if err != nil {
		panic(err)
	}
	return nil
}

func unzip(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer func(){
	  if err := reader.Close(); err != nil {
	    panic(err)
	  }
	}()

	os.MkdirAll(destination, 0755)

	for _, file := range reader.File {
		err := extractAndWriteFile(file, destination)
		if err != nil {
			return err
		}
	}

	return nil

}

func extractAndWriteFile(file *zip.File, destination string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer func(){
	  if err := rc.Close(); err != nil {
	    panic(err)
	  }
	}()

	path := filepath.Join(destination, file.Name)
	if file.FileInfo().IsDir() {
		os.MkdirAll(path, file.Mode())
	} else {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer closeFile(file)

		_, err = io.Copy(file, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func joinFiles(sourceBaseName string) error {
	files := []string{}
	err := filepath.Walk(sourceBaseName+"\\", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer func(fileName string) {
			if err := file.Close(); err != nil {
				panic(err)
			}
			deleteFile(fileName)
		}(fileName)

		buf := make([]byte, groupFileSize)
		for {
			n, err := file.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal("Error reading file: ", err)
				return err
			}

			if n == 0 {
				break
			}

			writeFile(sourceBaseName+"\\"+sourceBaseName, buf[:n])
		}
	}

	return nil
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		panic(err)
	}
}
