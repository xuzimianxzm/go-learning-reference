package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
)

// Wordbook
/**
网易单词本导出 xml 格式例子
<wordbook>
    <item>
        <word>engineering</word>
        <trans><![CDATA[n. 工程，工程学 v. 设计；管理（engineer的ing形式）；建造 ]]></trans>
        <phonetic><![CDATA[]]></phonetic>
        <tags></tags>
        <progress>-1</progress>
    </item>
</wordbook>
*/
type Wordbook struct {
	XMLName xml.Name `xml:"wordbook"` //标签上的标签名
	Items   []Item   `xml:"item"`
}
type Item struct {
	Word     string `xml:"word"`
	Trans    string `xml:"trans"`
	Tags     string `xml:"tags"`
	Progress int    `xml:"progress"`
}

func main() {
	content, err := ioutil.ReadFile("./单词本导出-2021-08-24.xml")
	if err != nil {
		fmt.Printf("读取文件失败：%v", err)
		return
	}

	var result Wordbook
	if err = xml.Unmarshal(content, &result); err != nil {
		fmt.Printf("解析xml文件失败：%v", err)
	}

	var words string
	for _, item := range result.Items {
		words = words + strings.TrimSpace(item.Word) + "\r"
	}

	ioutil.WriteFile("单词本导出-2021-08-24.txt", []byte(words), 0644)
}
