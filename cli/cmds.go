package cli

import (
	"flag"
	"fmt"
	"give-me-genshin-gacha/database"
	"give-me-genshin-gacha/gacha"
	"give-me-genshin-gacha/models"
	"os"
)

type Action = func(*flag.FlagSet) error
type FlagSet struct {
	*flag.FlagSet
	action Action
}

func (f *FlagSet) Run() error {
	return f.action(f.FlagSet)
}

func NewFlagSet(name, usage string, action Action) FlagSet {
	f := FlagSet{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
		action:  action,
	}
	f.Usage = func() {
		fmt.Fprintf(f.Output(), "%s : %s\n", name, usage)
		f.PrintDefaults()
	}
	return f
}

type Cli struct {
	commands map[string]FlagSet
}

func (c *Cli) Run(args []string) error {
	if cmd, ok := c.commands[args[0]]; ok {
		return cmd.Run()
	}
	fmt.Printf("unknow command [%s]\n", args[0])
	c.Usage()
	return nil
}
func (c *Cli) Usage() {
	for _, cmd := range c.commands {
		cmd.Usage()
	}
}

func NewCli() Cli {
	c := Cli{
		commands: map[string]FlagSet{},
	}
	cmdGameDir :=
		NewFlagSet("gamedir", "show game data dir", GetGameDir)
	c.commands[cmdGameDir.Name()] = cmdGameDir
	return c
}
func ShowErr(e error) {
	fmt.Printf("Error: %s", e)
	os.Exit(1)
}
func GetGameDir(*flag.FlagSet) error {
	dir, err := gacha.GetGameDir()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return err
}

func GetGachaURL(gameDir string, useProxy bool) {
	// TODO: 使用代理服务器
	if gameDir == "" {
		dir, err := gacha.GetGameDir()
		if err != nil {
			ShowErr(err)
		}
		gameDir = dir
	}
	url, err := gacha.GetRawURL(gameDir)
	if err != nil {
		ShowErr(err)
	}
	fmt.Println(url)
	os.Exit(0)
}

func UpdateItem(lang string) {
	if lang == "" {
		lang = "zh-cn"
	}
	db, err := database.NewDB("./data.db")
	if err != nil {
		panic(err)
	}
	avatar, weapon, err := gacha.GetGameItem(lang)
	if err != nil {
		ShowErr(err)
	}
	items := make([]models.Item, 0, len(avatar)+len(weapon))
	names := make([]models.ItemName, 0, len(avatar)+len(weapon))
	for _, a := range avatar {
		items = append(items, models.Item{
			ItemID:   a.ItemID,
			RankType: a.RankType,
			ItemType: models.ItemTypeAvatar,
		})
		names = append(names, models.ItemName{
			ItemID: a.ItemID,
			Lang:   lang,
			Name:   a.Name,
		})
	}
	for _, w := range weapon {
		items = append(items, models.Item{
			ItemID:   w.ItemID,
			RankType: w.RankType,
			ItemType: models.ItemTypeWeapon,
		})
		names = append(names, models.ItemName{
			ItemID: w.ItemID,
			Lang:   lang,
			Name:   w.Name,
		})
	}
	addItems, err := db.Items.Put(items...)
	if err != nil {
		ShowErr(err)
	}
	addItemNames, err := db.ItemNames.Put(names...)
	if err != nil {
		ShowErr(err)
	}
	type Info struct {
		ID                   uint64
		Name                 string
		isNewItem, isNewName bool
	}
	fmt.Println("ID\tName\tNewItem\tNewName")
	info := make(map[uint64]Info)
	for _, item := range addItems {
		info[item.ItemID] = Info{
			ID:        item.ItemID,
			isNewItem: true,
		}
	}
	for _, name := range addItemNames {
		i, ok := info[name.ItemID]
		i.isNewName = true
		i.Name = name.Name
		if !ok {
			i.ID = name.ItemID
		}
	}
	for _, i := range info {
		fmt.Printf("%d\t%s\t%t\t%t", i.ID, i.Name, i.isNewItem, i.isNewName)
	}
	os.Exit(0)
}
