package deck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var deckIDPattern = regexp.MustCompile(`^deck-(\d{4,})$`)

//计算出候选编号
func(s *Service) claimDeckDir() (string,error){
	for i:=0;i<100;i++{
		next,err := s.nextDeckNumber()
		if err !=nil{
			return "",err
		}

		id := fmt.Sprintf("deck-%04d",next)
		err = os.Mkdir(filepath.Join(s.decksDir,id),0o755)
		//名称被占用则重试
		if os.IsExist(err){
			continue
		}
		if err !=nil{
			return "",err
		}
		return id,nil
	}
	return "",errors.New("分配deck id 重试超限")
}

//扫描现有目录 返回最大编号+1
func(s *Service) nextDeckNumber()(int,error){
	entries,err := os.ReadDir(s.decksDir)
	if os.IsNotExist(err) {
		return 1,nil //一个文件都没有
	}

	if err!=nil{
		return 0,err
	}

	max := 0
	for _,e := range entries{
		m := deckIDPattern.FindStringSubmatch(e.Name())
		if m == nil{
			continue
		}
		if n,err := strconv.Atoi(m[1]);err == nil && n>max{
			max = n
		}
	}
	return max +1,nil
}
// atomicWriteFile 原子写：临时文件 + rename，防止写入一半的中间态被读到。
func atomicWriteFile(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
