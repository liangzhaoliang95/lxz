/**
 * @author  zhaoliang.liang
 * @date  2025/8/4 10:50
 */

package view

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/drivers/database_drivers"
	"lxz/internal/view/base"
)

type DatabaseDbTree struct {
	*BaseFlex
	app             *App
	tableChangeChan chan tableChangeSubscribe // 用于订阅表变化的通道
	// 数据库
	dbCfg        *config.DBConnection           // 数据库连接配置
	dbConn       database_drivers.IDatabaseConn // 数据库连接接口
	databaseList []string                       // 当前连接下的数据库列表
	tableList    []string                       // 当前数据库下的表列表
	selectDB     string                         // 当前选中的数据库
	selectTable  string                         // 当前选中的表
	// UI组件
	databaseUiTree *tview.TreeView // 用于显示数据库树
}

func (_this *DatabaseDbTree) selfFocus() {
	// 设置当前焦点为表格组件
	_this.app.UI.SetFocus(_this)
}

func (_this *DatabaseDbTree) Init(ctx context.Context) error {
	// 获取数据库连接配置
	// 初始化数据库连接
	iDatabaseConn, err := database_drivers.GetConnectOrInit(_this.dbCfg)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_this.dbConn = iDatabaseConn

	// 获取当前连接下的数据库列表
	dbList, err := _this.dbConn.GetDbList()
	if err != nil {
		return err
	}
	_this.databaseList = dbList
	slog.Info("databaseList", "dbName", _this.dbCfg.Name, "list", _this.databaseList)

	// 初始化tree view
	_this.databaseUiTree = tview.NewTreeView()
	_this.databaseUiTree.SetBorder(false)
	_this.databaseUiTree.SetTopLevel(1)
	_this.databaseUiTree.SetSelectedFunc(func(node *tview.TreeNode) {
		selectName := node.GetText()
		slog.Info("Selected Node", "selectName", selectName, "level", node.GetLevel())
		if len(node.GetChildren()) > 0 {
			// 如果节点已经有子节点，直接返回
			node.SetExpanded(!node.IsExpanded())
		} else {
			// 判断是否是顶级节点
			if node.GetLevel() == 1 {
				// 添加表节点
				tableList, err := _this.dbConn.GetTableList(selectName)
				if err != nil {
					slog.Error("Failed to get table list", "dbName", selectName, "error", err)
					return
				}
				for _, tableName := range tableList {
					tableNode := tview.NewTreeNode(tableName).
						SetColor(tview.Styles.SecondaryTextColor).
						SetSelectable(true)
					node.AddChild(tableNode)
				}
			} else {
				// 如果是表节点，渲染右侧
				parentNode := node.GetParentNode()
				dbName := parentNode.GetText()
				_this.selectDB = dbName

				tableName := node.GetText()
				_this.selectTable = tableName

				slog.Info("Selected Node is a table node", "dbName", dbName, "tableName", selectName)
				// 启动表视图 发送表变化订阅
				tableChangeChan := tableChangeSubscribe{
					dbName:    dbName,
					tableName: tableName,
				}
				_this.tableChangeChan <- tableChangeChan

			}

		}

	})

	_this.AddItem(_this.databaseUiTree, 0, 1, true)

	return nil
}

func (_this *DatabaseDbTree) Start() {
	slog.Info("DatabaseDbTree Start", "dbName", _this.dbCfg.Name)
	// 设置树视图的根节点
	rootNode := tview.NewTreeNode("-").
		SetColor(tview.Styles.PrimaryTextColor)
	_this.databaseUiTree.SetRoot(rootNode)
	_this.databaseUiTree.SetCurrentNode(rootNode)
	// 遍历数据库列表，添加到树视图中
	for _, dbName := range _this.databaseList {

		dbNode := tview.NewTreeNode(dbName).
			SetColor(tcell.ColorGold).
			SetSelectable(true)
		_this.databaseUiTree.GetRoot().AddChild(dbNode)

	}

}

func (_this *DatabaseDbTree) Stop() {

}

func NewDatabaseDbTree(
	a *App,
	dbCfg *config.DBConnection,
	tableChangeChan chan tableChangeSubscribe,
) *DatabaseDbTree {
	var name = dbCfg.Name
	lp := DatabaseDbTree{
		BaseFlex:        NewBaseFlex(name),
		app:             a,
		dbCfg:           dbCfg,
		tableChangeChan: tableChangeChan,
	}
	lp.SetBorder(true)
	lp.SetBorderColor(base.BoarderDefaultColor)
	return &lp
}
