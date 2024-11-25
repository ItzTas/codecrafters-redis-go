package main

//
// import "os"
//
// type RDBWriter struct{}
//
// func (r *RDBWriter) writeRDB(path string, metadata map[string]string) error {
// 	err := r.writeHeader(path)
// 	if err != nil {
// 		return err
// 	}
// 	err = r.writeMetadata(path, metadata)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (r *RDBWriter) writeMetadata(path string, metadata map[string]string) error {
// 	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	_, err = file.Write([]byte("redis-dump-path"))
// 	if err != nil {
// 		return err
// 	}
// 	_, err = file.Write([]byte(metadata["dir"] + metadata["filename"]))
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (r *RDBWriter) writeHeader(path string) error {
// 	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	_, err = file.Write([]byte("REDIS0011\n"))
// 	return err
// }
