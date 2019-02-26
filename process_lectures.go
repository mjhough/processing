package main

import(
  "fmt"
  "flag"
  "io/ioutil"
  "os"
  "bufio"
)

func main() {
  readDir := flag.String("rdir", "", "directory to read from")
  checkFile := flag.String("cfile", "", "file to check and write to")
  flag.Parse()

  /*
  Get list of files in directory. From there we can check
  whether or not these files have been processed already or not.
  */
  files, err := listFiles(*readDir)

  if err != nil {
    fmt.Println(err)
    return
  }

  /*
  Check what files need to be synced by seeing if there exist
  any discrepancies between the check file and the current list
  of files.
  */
  filesToSync, err := checkFilesToSync(files, *checkFile)

  if err != nil {
    fmt.Println(err)
    return
  }

  /*
  Sync all files that need to be synced.
  This means SCPing the files across to Euler from Gauss.
  There they will be processed depending on their file name.
  */
  syncFiles(filesToSync)

  /*
  Remember the files in the directory at last run
  by saving a list of them to a text file for future
  comparison.
  */
  rememberFiles(files, *checkFile)
}

func checkFilesToSync(files []os.FileInfo, checkFile string) ([]string, error) {
  cFile, err := os.Open(checkFile)
  if err != nil {
    return nil, err
    cFile.Close()
  }

  var lines []string
  scanner := bufio.NewScanner(cFile)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }

  if scanner.Err() != nil {
    return nil, scanner.Err()
    cFile.Close()
  }

  var fileNames []string
  for _, fi := range files {
    fileNames = append(fileNames, fi.Name())
  }

  return difference(fileNames, lines), nil
}

func difference(a, b []string) []string {
  mb := map[string]bool{}
  for _, x := range b {
    mb[x] = true
  }
  ab := []string{}
  for _, x := range a {
    if _, ok := mb[x]; !ok {
      ab = append(ab, x)
    }
  }
  return ab
}

func listFiles(readDir string) ([]os.FileInfo, error) {
  files, err := ioutil.ReadDir(readDir)
  if err != nil {
    return nil, fmt.Errorf("Error reading directory", err)
  }
  return files, nil
}

func rememberFiles(files []os.FileInfo, checkFile string) {
  f, err := os.Create(checkFile)
  if err != nil {
    fmt.Println(err)
    return
  }

  for _, fi := range files {
    _, err := fmt.Fprintln(f, fi.Name())
    if err != nil {
      fmt.Println(err)
      f.Close()
      return
    }
  }

  fmt.Println("Bytes written successfully")
  err = f.Close()
  if err != nil {
    fmt.Println(err)
    return
  }
}

func syncFiles(filesToSync []string) {
  // TODO: implement
}
