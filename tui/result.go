package tui

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/deckarep/golang-set"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResultView struct {
	*tview.Table
	ItemArray []Item
	Keys      []interface{}
}

type Item struct {
	Item map[string]interface{}
}

type Items struct {
	Items []map[string]interface{}
}

func NewResultView() *ResultView {
	rv := &ResultView{
		Table: tview.NewTable().Select(0, 0).SetSelectable(true, true),
	}
	rv.SetBorder(true).SetTitle("Results").SetTitleAlign(tview.AlignLeft)
	return rv
}

func (rv *ResultView) UpdateView(t *Tui) {
	table := rv.Clear()

	err := rv.RunCmd(t.QueryView.Query)
	if err != nil {
		messages := strings.Split(err.Error(), ": ")
		for i, msg := range messages {
			table.SetCell(i, 0, &tview.TableCell{
				Text:            msg,
				NotSelectable:   true,
				Align:           tview.AlignLeft,
				Color:           tcell.ColorRed,
				BackgroundColor: tcell.ColorDefault,
			})
		}
		return
	}
	rv.DrawResults()
}

func (rv *ResultView) RunCmd(sql string) error {
	buf := bytes.Buffer{}
	cmd := exec.Command("aws", "dynamodb", "execute-statement", "--statement", sql)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return errors.New(buf.String())
	}

	jsonStr := []byte(buf.String())
	var items Items
	if err := json.Unmarshal(jsonStr, &items); err != nil {
		return err
	}

	var list []Item
	for _, m := range items.Items {
		values := make(map[string]*dynamodb.AttributeValue)
		for k, mp := range m {
			var ddbAttr dynamodb.AttributeValue
			bytes, err := json.Marshal(mp)
			if err != nil {
				return err
			}
			json.Unmarshal(bytes, &ddbAttr)
			values[k] = &ddbAttr
		}
		var item Item
		dynamodbattribute.UnmarshalMap(values, &item.Item)
		list = append(list, item)
	}

	keys := mapset.NewSet()
	for _, item := range list {
		for k := range item.Item {
			keys.Add(k)
		}
	}
	keyArray := keys.ToSlice()
	sort.Slice(keyArray, func(i, j int) bool { return keyArray[i].(string) < keyArray[j].(string) })
	rv.Keys = keyArray

	rv.ItemArray = list
	return nil
}

func (rv *ResultView) DrawResults() {
	t := rv.Clear()
	c := 0
	for i, h := range rv.Keys {
		t.SetCell(0, c+i, &tview.TableCell{
			Text:            h.(string),
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorYellow,
			BackgroundColor: tcell.ColorDefault,
		})
	}

	for i, item := range rv.ItemArray {
		c := 0
		for j, key := range rv.Keys {
			if item.Item[key.(string)] == nil {
				t.SetCell(i+1, c+j, tview.NewTableCell(""))
			} else {
				json, err := json.Marshal(item.Item[key.(string)])
				if err != nil {
					return
				}
				t.SetCell(i+1, c+j, tview.NewTableCell(string(json)))
			}
		}
	}
}
