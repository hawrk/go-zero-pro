// Package logic
/*
 Author: hawrkchen
 Date: 2022/10/24 20:33
 Desc:
*/
package logic

import (
	"account-auth/account-auth-server/models"
	"encoding/json"
	"errors"
)

type UserMenuList struct {
	List string `json:"list"`
}

// MenuList 权限菜单列表
type MenuList struct {
	List []FirstMenu `json:"list"`
}

// FirstMenu 一级菜单列表
type FirstMenu struct {
	OneLevelName string       `json:"name"`
	Auth         int          `json:"auth"` //权限 1-拥有权限
	SecondM      []SecondMenu `json:"children"`
	Cmpt         []CmptAuth   `json:"cmpt"`
}

// SecondMenu 二级菜单列表
type SecondMenu struct {
	TwoLevelName string      `json:"name"`
	Auth         int         `json:"auth"` //权限 1-拥有权限
	ThirdM       []ThirdMenu `json:"children"`
	Cmpt         []CmptAuth  `json:"cmpt"`
}

// ThirdMenu 三级菜单列表
type ThirdMenu struct {
	ThirdLevelName string     `json:"name"`
	Auth           int        `json:"auth"` //权限 1-拥有权限
	Cmpt           []CmptAuth `json:"cmpt"`
}

// CmptAuth 权限组件
type CmptAuth struct {
	Name     string `json:"name"`
	Auth     int    `json:"auth"`
	CmptType int    `json:"type"` // 控件类型：1-查询2-新增3-上传4-导出列表5-下载报告6-删除
}

func BuildMenuList(userType int32, out []models.TbAuthMenu) (string, error) {
	if len(out) <= 0 {
		return "", errors.New("empty auth info")
	}
	//var ret MenuList
	var fms []FirstMenu
	// 1。 先关联父菜单和子菜单之间的关系
	var OneLevel []int // 一级菜单列表
	var m2 []int       // 二级菜单列表
	var m3 []int       // 三级菜单列表
	TwoLevel := make(map[int][]int)
	ThirdLevel := make(map[int][]int)
	ComponentMap := make(map[int][]int) // 控件list, key： 菜单id， value -> 控件列表
	CmptTypeMap := make(map[int]int)    // 控件类型， key: 控件ID（menuid) , value-> 控件类型

	menuId2Name := make(map[int]string) // key -> menuid, value -> menuname
	menuId2Auth := make(map[int]int)    // key -> menuid, value -> 默认权限
	for _, v := range out {
		// 先把一二三级菜单找出来
		if v.MenuType == 1 {
			OneLevel = append(OneLevel, v.MenuId)
			//m1[v.MenuId] = v.MenuName
		} else if v.MenuType == 2 {
			m2 = append(m2, v.MenuId)
			//m2[v.MenuId] = v.MenuName
		} else if v.MenuType == 3 {
			//m3[v.MenuId] = v.MenuName
			m3 = append(m3, v.MenuId)
		} else if v.MenuType == 11 {
			ComponentMap[v.ParManuId] = append(ComponentMap[v.ParManuId], v.MenuId)
			CmptTypeMap[v.MenuId] = v.CmptType
		}

		menuId2Name[v.MenuId] = v.MenuName
		if userType == 1 { //如果是超级管理员，则全部返回
			menuId2Auth[v.MenuId] = 1
		} else {
			menuId2Auth[v.MenuId] = v.AuthDefault
		}

	}
	//l.Logger.Infof("get ComponentMap:%+v", ComponentMap)
	//for k, v := range ComponentMap {
	//	l.Logger.Info("get key:", k, ", value:", v)
	//}
	// 一级菜单填充二级菜单
	for _, v := range OneLevel {
		var s2 []int
		for _, v1 := range out {
			if v == v1.ParManuId && v1.MenuType != 11 { //排除掉控件
				s2 = append(s2, v1.MenuId)
			}
		}
		TwoLevel[v] = s2
	}
	//l.Logger.Infof("get TwoLevel:%+v", TwoLevel)

	// 二级菜单填充三级菜单
	for _, v := range m2 {
		var s3 []int
		for _, v2 := range out {
			if v == v2.ParManuId && v2.MenuType != 11 {
				s3 = append(s3, v2.MenuId)
			}
		}
		ThirdLevel[v] = s3
	}
	//l.Logger.Infof("get ThirdLevel:%+v", ThirdLevel)

	// 2. 拼装Json串
	for _, v1 := range OneLevel {
		var fm FirstMenu
		fm.OneLevelName = GetManuNameByMenuId(menuId2Name, v1) //填充第一层
		fm.Auth = GetAuthByMenuId(menuId2Auth, v1)
		fm.Cmpt = GetCmptAuth(menuId2Auth, menuId2Name, CmptTypeMap, ComponentMap[v1])
		// 检查该层是否有二级菜单
		var sms []SecondMenu
		for _, v2 := range TwoLevel[v1] { // 有二级菜单
			var sm SecondMenu
			sm.TwoLevelName = GetManuNameByMenuId(menuId2Name, v2)
			sm.Auth = GetAuthByMenuId(menuId2Auth, v2)
			sm.Cmpt = GetCmptAuth(menuId2Auth, menuId2Name, CmptTypeMap, ComponentMap[v2])
			var tms []ThirdMenu
			for _, v3 := range ThirdLevel[v2] { // 有三级菜单
				var tm ThirdMenu
				tm.ThirdLevelName = GetManuNameByMenuId(menuId2Name, v3)
				tm.Auth = GetAuthByMenuId(menuId2Auth, v3)
				tm.Cmpt = GetCmptAuth(menuId2Auth, menuId2Name, CmptTypeMap, ComponentMap[v3])
				tms = append(tms, tm)
			}
			sm.ThirdM = tms
			sms = append(sms, sm)
		}
		fm.SecondM = sms
		fms = append(fms, fm)
		//ret.List = append(ret.List, fm)
	}
	//l.Logger.Infof("get ret:%+v", ret)
	o, err := json.Marshal(fms)
	if err != nil {
		return "", err
		//l.Logger.Error("marshal error:", err)
	}
	//l.Logger.Info("get json result:", string(o))

	//var other MenuList
	//if err := json.Unmarshal(o, &other); err != nil {
	//	l.Logger.Error("unmarshal error:", err)
	//}
	//l.Logger.Infof("get other:%+v", other)

	return string(o), nil
}

// GetManuNameByMenuId 根据菜单ID 映射出菜单名称
func GetManuNameByMenuId(m map[int]string, id int) string {
	if _, exist := m[id]; exist {
		return m[id]
	}
	return ""
}

// GetAuthByMenuId 根据菜单ID查其默认权限
func GetAuthByMenuId(m map[int]int, id int) int {
	if _, exist := m[id]; exist {
		return m[id]
	}
	return 0
}

// GetCmptAuth 根据传入的菜单ID，返回其所属于控件的权限
func GetCmptAuth(m map[int]int, m2 map[int]string, m3 map[int]int, ids []int) []CmptAuth {
	var out []CmptAuth
	for _, v := range ids {
		c := CmptAuth{
			Name:     GetManuNameByMenuId(m2, v),
			Auth:     GetAuthByMenuId(m, v),
			CmptType: GetCmptType(m3, v),
		}
		out = append(out, c)
	}
	return out
}

func GetCmptType(m map[int]int, id int) int {
	if _, exist := m[id]; exist {
		return m[id]
	}
	return 0
}
