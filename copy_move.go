package copy_go

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// CopyFile is the simplest copy, copies file from src to dst
func CopyFile(src, dst string) (err error) {
	var input, output *os.File

	input, err = os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("src file does not exist")
		}
		return err
	}
	defer fclose(input, &err)

	output, err = os.Create(dst)
	if err != nil {
		return err
	}
	defer fclose(output, &err)

	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	return output.Sync()
}

// Move file from src to dst
func Move(src, dst string) (err error) {
	var input, output *os.File

	input, err = os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("src file does not exist")
		}
		return err
	}
	defer fclose(input, &err)

	output, err = os.Create(dst)
	if err != nil {
		return err
	}
	defer fclose(output, &err)

	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	if err = output.Sync(); err != nil {
		return err
	}
	return os.Remove(src)
}
