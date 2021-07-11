package core

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

type Html struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ElType    string `json:"elType"`
	TagName   string `json:"tagName"`
	ClassName string `json:"className"`
	Xpath     string `json:"xpath"`
}

// JSCode 获取表单下所有的输入框
const JSCode = `data()

				function data() {
    				let form = document.forms[%d]
    				let array = []
    				for (let i = 0; i < form.elements.length; i++) {
    				    let ee = form.elements[i];
    				    if ("INPUT" === ee.tagName || "SELECT" === ee.tagName) {
    				        array.push({ 'tagName': ee.tagName, 'name': ee.name, 'elType': ee.type, 'id': ee.id, 'className': ee.className, 'xpath': getPathTo(ee)
    				        })
    				    }
    				} return array }

				function getPathTo(element) {
					if (element.id !== '')
						return 'id("' + element.id + '")';
					if (element === document.body)
						return element.tagName;
					let index = 0;
					let siblings = element.parentNode.childNodes;
					for (let i = 0; i < siblings.length; i++) {
						let sibling = siblings[i];
						if (sibling === element)
							return getPathTo(element.parentNode) + '/' + element.tagName + '[' + (index + 1) + ']';
						if (sibling.nodeType === 1 && sibling.tagName === element.tagName)
							index++;
					}
				}`

// inputMap 全局input接口类型
var inputMap = map[string]Input{
	"text":     InputText{},
	"password": InputPassword{},
	// "email":      InputEmail{},
	"date":       InputDate{},
	"radio":      InputRadio{},
	"checkbox":   InputCheckbox{},
	"select-one": InputSelect{},
	"submit":     InputSubmit{},
}

func Form(ctx context.Context) {

	var formLength int

	chromedp.Run(ctx,
		chromedp.Evaluate("document.forms.length", &formLength))

	for i := 0; i < formLength; i++ {
		var res []Html
		chromedp.Run(ctx, chromedp.Evaluate(fmt.Sprintf(JSCode, i), &res))
		for _, re := range res {
			fmt.Println("打印ElType", re.ElType)
			if input, ok := inputMap[re.ElType]; ok {
				input.handler(ctx, re)
			}
		}
	}
}

// Input 输入接口
type Input interface {
	// 执行输入接口的方法
	handler(context.Context, Html)
}

type InputText struct {
}

func (text InputText) handler(ctx context.Context, html Html) {
	chromedp.Run(ctx, chromedp.SendKeys(html.Xpath, "mengjiaheng", chromedp.AtLeast(0)))
	// utils.Run(ctx, chromedp.SendKeys(html.Xpath, utils.RandomName(), chromedp.AtLeast(0)))
}

type InputPassword struct {
}

func (password InputPassword) handler(ctx context.Context, html Html) {
	chromedp.Run(ctx, chromedp.SendKeys(html.Xpath, "123456", chromedp.AtLeast(0)))
	// utils.Run(ctx, chromedp.SendKeys(html.Xpath, utils.RandomPassword(), chromedp.AtLeast(0)))
}

type InputEmail struct {
}

// func (email InputEmail) handler(ctx context.Context, html Html) {
// 	utils.Run(ctx, chromedp.SendKeys(html.Xpath, utils.RandomEmail(), chromedp.AtLeast(0)))
// }

type InputDate struct {
}

func (email InputDate) handler(ctx context.Context, html Html) {
	chromedp.Run(ctx, chromedp.SendKeys(html.Xpath, "测试1", chromedp.AtLeast(0)))

	// utils.Run(ctx, chromedp.SendKeys(html.Xpath, utils.RandomDate(), chromedp.AtLeast(0)))
}

type InputSubmit struct {
}

func (submit InputSubmit) handler(ctx context.Context, html Html) {
	chromedp.Run(ctx, chromedp.Submit(html.Xpath, chromedp.AtLeast(0)), chromedp.Stop())

	// utils.Run(ctx, chromedp.Submit(html.Xpath, chromedp.AtLeast(0)), chromedp.Stop())
}

type InputRadio struct {
}

func (radio InputRadio) handler(ctx context.Context, html Html) {
	chromedp.Run(ctx, chromedp.Click(html.Xpath, chromedp.AtLeast(0)))
	// utils.Run(ctx, chromedp.Click(html.Xpath, chromedp.AtLeast(0)))
}

type InputSelect struct {
}

//不知名函数
func IsCdpNode(arr []*cdp.Node) bool {
	if arr == nil {
		return true
	}
	return false
}
func IsNotIsCdpNode(arr []*cdp.Node) bool {
	return !IsCdpNode(arr)
}

func (selects InputSelect) handler(ctx context.Context, html Html) {

	chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		_ = chromedp.Click(html.Xpath, chromedp.AtLeast(0)).Do(ctx)
		var projects []*cdp.Node
		_ = chromedp.Nodes(html.Xpath, &projects, chromedp.AtLeast(0)).Do(ctx)
		if IsNotIsCdpNode(projects) {
			for _, child := range projects[0].Children {
				optgroup, _ := child.Attribute("optgroup")
				if optgroup == "" {
					_ = chromedp.SendKeys(html.Xpath, kb.ArrowDown).Do(ctx)
					break
				}
			}
		}
		return nil
	}))
}

type InputCheckbox struct {
}

func (checkbox InputCheckbox) handler(ctx context.Context, html Html) {
	chromedp.Run(ctx, chromedp.Click(html.Xpath, chromedp.AtLeast(0)))
}
