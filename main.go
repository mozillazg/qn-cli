/*
七牛本地上传客户端
$ qn_cli --help
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/resumable/io"
	"github.com/qiniu/api/rs"
)

// 生成上传 token
func genToken(bucket string) string {
	policy := rs.PutPolicy{
		Scope:      bucket,
		ReturnBody: `{"bucket": $(bucket),"key": $(key)}`,
	}
	return policy.Token(nil)
}

// 上传本地文件
func uploadFile(localFile, key, uptoken string) (ret io.PutRet, err error) {
	var extra = &io.PutExtra{}

	err = io.PutFile(nil, &ret, uptoken, key, localFile, extra)
	if err != nil {
		fmt.Println("Upload file failed:", err)
		return
	}

	return
}

func main() {
	// 文件路径
	file := flag.String("f", "--file", "local file")
	// 保存名称
	saveName := flag.String("s", "--save", "save name")
	flag.Parse()

	bucketName := os.Getenv("QINIU_BUCKET_NAME")
	bucketURL := os.Getenv("QINIU_BUCKET_URL")
	accessKey := os.Getenv("QINIU_ACCESS_KEY")
	secretKey := os.Getenv("QINIU_SECRET_KEY")
	key := *saveName

	if *file == "" {
		flag.PrintDefaults()
		fmt.Println("need file: -f or --file")
		return
	}

	// 配置 accesskey, secretkey
	conf.ACCESS_KEY = accessKey
	conf.SECRET_KEY = secretKey
	// 生成上传 token
	uptoken := genToken(bucketName)

	ret, err := uploadFile(*file, key, uptoken)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(bucketURL + ret.Key)
	}
}
