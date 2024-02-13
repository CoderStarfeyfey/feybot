/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/23
 **/
package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CoderStarfeyfey/feybot/internal"
	"github.com/CoderStarfeyfey/feybot/internal/utils"
	"net/http"
	"path"
)

const baseURL = "https://leetcode-cn.com"

func GetDailyQuestionTitleSlug() (string, error) {
	query := `{"operationName":"questionOfToday","variables":{},"query":"query questionOfToday { todayRecord { question { questionFrontendId questionTitleSlug __typename } lastSubmission { id __typename } date userStatus __typename }}"}`

	resp, err := http.Post(baseURL+"/graphql", "application/json", bytes.NewBufferString(query))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			TodayRecord []struct {
				Question struct {
					QuestionTitleSlug string `json:"questionTitleSlug"`
				} `json:"question"`
			} `json:"todayRecord"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Data.TodayRecord) > 0 {
		return result.Data.TodayRecord[0].Question.QuestionTitleSlug, nil
	}
	return "", fmt.Errorf("no daily question found")
}

func GetQuestionData(titleSlug string) (*Question, error) {
	query := fmt.Sprintf(`{"operationName":"questionData","variables":{"titleSlug":"%s"},"query":"query questionData($titleSlug: String!) { question(titleSlug: $titleSlug) { questionId questionFrontendId boundTopicId title titleSlug content translatedTitle translatedContent isPaidOnly difficulty likes dislikes isLiked similarQuestions contributors { username profileUrl avatarUrl __typename } langToValidPlayground topicTags { name slug translatedName __typename } companyTagStats codeSnippets { lang langSlug code __typename } stats hints solution { id canSeeDetail __typename } status sampleTestCase metaData judgerAvailable judgeType mysqlSchemas enableRunCode envInfo book { id bookName pressName source shortDescription fullDescription bookImgUrl pressImgUrl productUrl __typename } isSubscribed isDailyQuestion dailyRecordStatus editorType ugcQuestionId style __typename }}"}`, titleSlug)

	resp, err := http.Post(baseURL+"/graphql", "application/json", bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Question *Question `json:"question"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data.Question, nil
}

// Question 结构体用于解析问题的 JSON 数据
type Question struct {
	QuestionFrontendId string `json:"questionFrontendId"`
	TranslatedTitle    string `json:"translatedTitle"`
	Difficulty         string `json:"difficulty"`
	CodeURL            string `json:"CodeURL"`
}

func init() {
	internal.PluginMap["LeetcodeReply"] = LeetcodeReply
}

func LeetcodeReply(req *internal.RequestStruct) (*internal.ReplyStruct, error) {
	questionTitleSlug, err := GetDailyQuestionTitleSlug()
	everyURL := path.Join(baseURL, questionTitleSlug)
	if err != nil {
		utils.FeyLog.Errorf("Error fetching daily question:", err)
		return nil, err
	}
	// 获取每日一题的所有信息
	questionData, err := GetQuestionData(questionTitleSlug)
	questionData.CodeURL = everyURL
	if err != nil {
		utils.FeyLog.Errorf("Error fetching question data:", err)
		return nil, err
	}
	formattedString := fmt.Sprintf("题号: %s\n题名(中文): %s\n难度级别: %s\n题目链接: %s\n",
		questionData.QuestionFrontendId, questionData.TranslatedTitle, questionData.Difficulty, questionData.CodeURL)
	utils.FeyLog.Debugf("Get the leetcode message:%v", formattedString)
	return &internal.ReplyStruct{
		ReType: 1,
		ReText: formattedString,
	}, nil
}

func Test() {
	questionTitleSlug, err := GetDailyQuestionTitleSlug()
	everyURL := path.Join(baseURL, questionTitleSlug)
	if err != nil {
		fmt.Println("Error fetching daily question:", err)
		return
	}

	// 获取每日一题的所有信息
	questionData, err := GetQuestionData(questionTitleSlug)
	questionData.CodeURL = everyURL
	if err != nil {
		fmt.Println("Error fetching question data:", err)
		return
	}
	fmt.Println("链接", questionData.CodeURL)
	fmt.Println("题号:", questionData.QuestionFrontendId)
	fmt.Println("题名（中文）:", questionData.TranslatedTitle)
	fmt.Println("难度级别:", questionData.Difficulty)
}
