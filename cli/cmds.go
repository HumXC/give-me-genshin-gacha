package cli

import (
	"flag"
	"fmt"
	"give-me-genshin-gacha/database"
	"give-me-genshin-gacha/gacha"
	"give-me-genshin-gacha/models"
	"strconv"
	"unicode/utf8"
)

type ActionMaker = func(*flag.FlagSet) Action
type Action = func() error
type FlagSet struct {
	*flag.FlagSet
	action Action
}

func (f *FlagSet) Run(args []string) error {
	err := f.Parse(args)
	if err != nil {
		return err
	}
	return f.action()
}

func NewFlagSet(name, usage string, actionMaker ActionMaker) FlagSet {
	set := flag.NewFlagSet(name, flag.ExitOnError)
	f := FlagSet{
		FlagSet: set,
		action:  actionMaker(set),
	}
	f.Usage = func() {
		fmt.Fprintf(f.Output(), "%s: %s\n", name, usage)
		f.PrintDefaults()
	}
	return f
}

type Cli struct {
	commands map[string]FlagSet
}

func (c *Cli) Run(args []string) {
	if cmd, ok := c.commands[args[0]]; ok {
		err := cmd.Run(args[1:])
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	fmt.Printf("Unknow command [%s]\n", args[0])
	c.Usage()
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
	cmdD :=
		NewFlagSet("D", "查找游戏的数据文件夹", GetGameDir)
	cmdG := NewFlagSet("G", "获取祈愿链接", GetGachaURL)
	cmdI := NewFlagSet("I", "更新物品信息", UpdateItem)
	c.commands[cmdD.Name()] = cmdD
	c.commands[cmdG.Name()] = cmdG
	c.commands[cmdI.Name()] = cmdI

	return c
}

func GetGameDir(*flag.FlagSet) Action {
	return func() error {
		dir, err := gacha.GetGameDir()
		if err != nil {
			return err
		}
		fmt.Println(dir)
		return err
	}
}

func GetGachaURL(flag *flag.FlagSet) Action {
	gameDir := flag.String("d", "", "指定游戏的目录，如果不指定则自动获取")
	// useProxy := flag.Bool("p", false, "使用代理服务器的方式获取链接")
	return func() error {
		var dir string
		// TODO: 使用代理服务器
		if *gameDir == "" {
			_dir, err := gacha.GetGameDir()
			if err != nil {
				return err
			}
			dir = _dir
		}
		url, err := gacha.GetRawURL(dir)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	}
}

func UpdateItem(flag *flag.FlagSet) Action {
	lang := flag.String("l", "zh-cn", "物品名称所使用的语言")
	dbFile := flag.String("d", "./data.db", "数据库的位置")
	return func() error {
		db, err := database.NewDB(*dbFile)
		if err != nil {
			return err
		}
		avatar, weapon, err := gacha.GetGameItem(*lang)
		if err != nil {
			return err
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
				Lang:   *lang,
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
				Lang:   *lang,
				Name:   w.Name,
			})
		}
		addItems, err := db.Items.Put(items...)
		if err != nil {
			return err
		}
		addItemNames, err := db.ItemNames.Put(names...)
		if err != nil {
			return err
		}
		type Info struct {
			ID                   string
			Name                 string
			isNewItem, isNewName bool
		}
		info := make(map[uint64]Info)
		for _, item := range addItems {
			info[item.ItemID] = Info{
				ID:        strconv.FormatUint(item.ItemID, 10),
				isNewItem: true,
			}
		}
		for _, name := range addItemNames {
			i, ok := info[name.ItemID]
			i.isNewName = true
			i.Name = name.Name
			if !ok {
				i.ID = strconv.FormatUint(name.ItemID, 10)
			}
			info[name.ItemID] = i
		}
		var lenID, lenName int
		for _, i := range info {
			_lenID := utf8.RuneCountInString(i.ID)
			if _lenID > lenID {
				lenID = _lenID
			}
			_lenName := utf8.RuneCountInString(i.Name)
			if _lenName > lenName {
				lenName = _lenName
			}
		}
		index := 0
		for _, i := range info {
			index++
			newItemFlag := " "
			newNameFlag := " "
			if i.isNewItem {
				newItemFlag = "*"
			}
			if i.isNewName {
				newNameFlag = "*"
			}
			fmt.Printf("[%d]\t %s%-"+strconv.Itoa(lenID)+"s\t %s%-"+strconv.Itoa(lenName)+"s\n", index, newItemFlag, i.ID, newNameFlag, i.Name)
		}
		return nil
	}
}
