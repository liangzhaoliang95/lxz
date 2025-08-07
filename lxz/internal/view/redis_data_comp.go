/**
 * @author  zhaoliang.liang
 * @date  2025/8/4 10:50
 */

package view

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"lxz/internal/config"
	"lxz/internal/helper"
	"lxz/internal/model"
	"lxz/internal/redis_drivers"
	"lxz/internal/ui/dialog"
	"lxz/internal/view/base"
)

type RedisDataComponent struct {
	*BaseFlex
	app             *App
	dbNum           int // 数据库名称
	redisConnConfig *config.RedisConnConfig
	rdbClient       *redis_drivers.RedisClient
	// 数据
	redisData *model.RedisData // Redis数据模型

	// ui组件
	filterFlex  *tview.Flex       // 用于布局过滤条件输入框和标签
	filterLabel *tview.TextView   // 用于显示过滤条件标签
	filterInput *tview.InputField // 用于输入过滤条件

	keyGroupFlex *tview.Flex     // 用于布局键分组树和键信息显示
	keyGroupTree *tview.TreeView // 用于显示键分组树

	keyInfoFlex *tview.Flex     // 用于显示键的详细信息
	keyName     *tview.TextView // 键名称
	KeyType     *tview.TextView // 键类型
	keyTTL      *tview.TextView // 键的TTL
	KeyValue    *tview.TextView // 键的内容
}

// refreshData
func (_this *RedisDataComponent) refreshData() error {
	// 刷新数据
	slog.Info("Refreshing Redis data", "dbNum", _this.dbNum)
	err := _this.refreshRedisData()
	if err != nil {
		slog.Error("Failed to refresh Redis data", "error", err)
		_this.app.UI.Flash().Err(fmt.Errorf("failed to refresh Redis data: %w", err))
		return err
	}
	// 刷新树形数据
	_this.SetTreeData()
	// 焦点切换到键分组树
	_this.focusKeyGroupTree()
	return nil
}

// newKey
func (_this *RedisDataComponent) newKey() {
	dialog.ShowCreateUpdateRedisData(&config.Dialog{}, _this.app.Content.Pages, &dialog.CreateUpdateRedisDataOpts{
		Title:   "New Key",
		Message: "",
		Data: &redis_drivers.RedisData{
			KetType: "string",
			KeyTTL:  -1,
		},
		Ack: func(data *redis_drivers.RedisData) bool {
			// 创建新键
			slog.Info("Creating new key", "data", data)
			if err := _this.rdbClient.CreateKeyData(data); err != nil {
				slog.Error("Failed to create new key", "error", err)
				_this.app.UI.Flash().Err(fmt.Errorf("failed to create new key: %w", err))
				return false
			}
			// 刷新数据
			err := _this.refreshRedisData()
			if err != nil {
				return true
			}
			// 刷新树形数据
			_this.SetTreeData()
			// 焦点切换到键分组树
			_this.focusKeyGroupTree()
			return true
		},
		Cancel: func() {
			// 取消操作，焦点切换到键分组树
			_this.focusKeyGroupTree()
		},
	})
}

// editKey
func (_this *RedisDataComponent) editKey() {
	// 获取当前选中的键
	selectedNode := _this.keyGroupTree.GetCurrentNode()
	slog.Info("editKey", "selectedNode", selectedNode)
	key := selectedNode.GetReference().(string)
	keyData, err := _this.rdbClient.GetKeyData(key)
	if err != nil {
		slog.Error("Failed to get key data", "key", key, "error", err)
		_this.app.UI.Flash().Err(fmt.Errorf("failed to get key data: %w", err))
		return
	}

	dialog.ShowCreateUpdateRedisData(&config.Dialog{}, _this.app.Content.Pages, &dialog.CreateUpdateRedisDataOpts{
		Title:   "Edit Key",
		Message: "",
		Data:    keyData,
		Ack: func(data *redis_drivers.RedisData) bool {
			// 更新键数据
			slog.Info("Updating key data", "key", key, "data", data)
			if err := _this.rdbClient.EditKeyData(key, data); err != nil {
				slog.Error("Failed to edit key data", "key", key, "error", err)
				_this.app.UI.Flash().Err(fmt.Errorf("failed to edit key data: %w", err))
				return false
			}
			// 刷新数据
			err = _this.refreshRedisData()
			if err != nil {
				return true
			}
			// 刷新树形数据
			_this.SetTreeData()
			// 焦点切换到键分组树
			_this.focusKeyGroupTree()
			return true
		},
		Cancel: func() {},
	})
}

func (_this *RedisDataComponent) deleteKey() {
	// 获取当前选中的键
	selectedNode := _this.keyGroupTree.GetCurrentNode()
	slog.Info("deleteKey", "selectedNode", selectedNode)
	key := selectedNode.GetReference().(string)
	dialog.ShowDelete(&config.Dialog{}, _this.app.Content.Pages, key, func(force bool) {
		_this.rdbClient.DeleteKey(key)
		parent := selectedNode.GetParentNode()
		parent.RemoveChild(selectedNode)

		_this.focusKeyGroupTree()
	}, func() {

		_this.focusKeyGroupTree()
	})

}

func (_this *RedisDataComponent) focusSearch() {
	// 设置当前焦点为搜索框
	_this.app.UI.SetFocus(_this.filterFlex)
	_this.filterFlex.SetBorderColor(base.ActiveBorderColor)
}

func (_this *RedisDataComponent) focusKeyGroupTree() {
	_this.app.UI.SetFocus(_this.keyGroupTree)
}

func (_this *RedisDataComponent) focusKeyInfoFlex() {
	_this.app.UI.SetFocus(_this.KeyValue)
}

func (_this *RedisDataComponent) refreshRedisData() error {
	// 初始化数据库连接
	iRedisConn, err := redis_drivers.GetConnectOrInit(_this.redisConnConfig, _this.dbNum)

	if err != nil {
		return fmt.Errorf("failed to get db connection: %w", err)
	}
	_this.rdbClient = iRedisConn

	// 初始化数据
	records, err := _this.rdbClient.GetRecords(
		"",
	)
	if err != nil {
		_this.app.UI.Flash().
			Err(fmt.Errorf("failed to get records for db %d: %w", _this.dbNum, err))
		return fmt.Errorf("failed to get records for db %d: %w", _this.dbNum, err)
	}
	// 对key分组
	_this.redisData = model.NewRedisData(records)
	return nil
}

func (_this *RedisDataComponent) Init(ctx context.Context) error {
	slog.Info("RedisDataComponent component init", "dbNum", _this.dbNum)

	err := _this.refreshRedisData()
	if err != nil {
		return err
	}
	// 初始化filterFlex
	_this.filterFlex = tview.NewFlex()
	_this.filterFlex.SetDirection(tview.FlexColumn)
	_this.filterFlex.SetBorder(true)
	_this.filterFlex.SetBorderPadding(0, 0, 1, 1)

	// 初始化filterLabel
	_this.filterLabel = tview.NewTextView()
	_this.filterLabel.SetText("Search: ")
	_this.filterLabel.SetTextAlign(tview.AlignCenter)
	_this.filterLabel.SetTextColor(tcell.ColorGreen)
	_this.filterLabel.SetBorderPadding(0, 0, 0, 0)
	_this.filterFlex.AddItem(_this.filterLabel, 8, 1, false)

	// 初始化filterInput
	_this.filterInput = tview.NewInputField()
	_this.filterInput.SetPlaceholder("Enter a clause to filter the results")
	_this.filterInput.SetFieldBackgroundColor(tcell.ColorBlack)
	_this.filterInput.SetFieldTextColor(tcell.ColorRed)
	_this.filterInput.SetFocusFunc(func() {
		_this.filterFlex.SetBorderColor(base.ActiveBorderColor)
	})
	_this.filterInput.SetBlurFunc(func() {
		_this.filterFlex.SetBorderColor(base.InactiveBorderColor)
	})
	_this.filterInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			whereClause := ""
			if _this.filterInput.GetText() != "" {
				whereClause = _this.filterInput.GetText()
			}
			slog.Info("filter input done", "where", whereClause)
			// 初始化表格数据
			keys, err := _this.rdbClient.GetRecords(
				whereClause,
			)
			_this.redisData = model.NewRedisData(keys)
			if err != nil {
				slog.Error(
					"Failed to get records for db",
					"dbNum",
					_this.dbNum,
					"error",
					err,
				)
				_this.app.UI.Flash().
					Err(fmt.Errorf("%w", err))

				_this.focusSearch()
			} else {
				_this.SetTreeData()

				// 焦点切换到表格
				_this.focusKeyGroupTree()
			}
		case tcell.KeyEscape:

		}
	})
	_this.filterFlex.AddItem(_this.filterInput, 0, 5, true)
	_this.AddItem(_this.filterFlex, 3, 0, false)

	// 初始化keyGroupFlex
	_this.keyGroupFlex = tview.NewFlex()
	_this.keyGroupFlex.SetDirection(tview.FlexColumn)
	_this.keyGroupFlex.SetBorder(false)
	_this.AddItem(_this.keyGroupFlex, 0, 1, true)

	// 初始化keyGroupTree
	_this.keyGroupTree = tview.NewTreeView()
	_this.keyGroupTree.SetBorder(true)
	_this.keyGroupTree.SetBorderPadding(0, 1, 1, 1)
	_this.keyGroupTree.SetTopLevel(1)
	_this.keyGroupTree.SetSelectedFunc(func(node *tview.TreeNode) {
		selectName := node.GetText()
		slog.Info(
			"redis key group selected Node",
			"selectName",
			selectName,
			"level",
			node.GetLevel(),
		)
		if len(node.GetChildren()) > 0 {
			// 如果节点已经有子节点，直接返回
			node.SetExpanded(!node.IsExpanded())
		} else {
			// 如果没有子节点，获取当前节点的完整路径
			fullPath := node.GetReference().(string)
			// key
			_this.keyName.SetText(fullPath)

			// 获取key值
			keyValue, err := _this.rdbClient.GetKeyValue(fullPath)
			if err != nil {
				slog.Error("Failed to get key value", "key", fullPath, "error", err)
				_this.app.UI.Flash().Err(fmt.Errorf("failed to get key value: %w", err))
				return
			}
			slog.Info(fmt.Sprintf("%s", helper.Prettify(keyValue)))
			_this.KeyValue.SetText(helper.Prettify(keyValue))
			// 获取key类型
			keyType, err := _this.rdbClient.GetKeyType(fullPath)
			if err != nil {
				slog.Error("Failed to get key type", "key", fullPath, "error", err)
				_this.app.UI.Flash().Err(fmt.Errorf("failed to get key type: %w", err))
				return
			}
			_this.KeyType.SetText(keyType)

			// 获取TTL
			ttl, err := _this.rdbClient.GetKeyTTL(fullPath)
			if err != nil {
				slog.Error("Failed to get key TTL", "key", fullPath, "error", err)
				_this.app.UI.Flash().Err(fmt.Errorf("failed to get key TTL: %w", err))
				return
			}
			_this.keyTTL.SetText(fmt.Sprintf("%d", ttl))

			slog.Info("Selected key info", "TTL", ttl, "key", fullPath, "type", keyType, "value", keyValue)
			_this.focusKeyInfoFlex()

		}
	})
	_this.keyGroupTree.SetFocusFunc(func() {
		// 设置当前焦点为键分组树
		_this.keyGroupTree.SetBorderColor(base.ActiveBorderColor)
	})
	_this.keyGroupTree.SetBlurFunc(func() {
		// 设置当前焦点为键分组树
		_this.keyGroupTree.SetBorderColor(base.InactiveBorderColor)
	})
	_this.keyGroupFlex.AddItem(_this.keyGroupTree, 0, 1, true)

	// 初始化keyInfoForm
	_this.keyInfoFlex = tview.NewFlex()
	_this.keyInfoFlex.SetBorder(false)
	_this.keyInfoFlex.SetDirection(tview.FlexRow)
	_this.keyGroupFlex.AddItem(_this.keyInfoFlex, 0, 1, true)

	// 初始化keyName
	_this.keyName = tview.NewTextView()
	_this.keyName.SetBorder(true)
	_this.keyName.SetTitle("KEY")
	_this.keyName.SetTitleAlign(tview.AlignCenter)
	_this.keyInfoFlex.AddItem(_this.keyName, 3, 1, false)

	// 初始化KeyType
	_this.KeyType = tview.NewTextView()
	_this.KeyType.SetBorder(true)
	_this.KeyType.SetTitle("TYPE")
	_this.KeyType.SetTitleAlign(tview.AlignCenter)
	_this.keyInfoFlex.AddItem(_this.KeyType, 3, 1, false)

	// 初始化keyTTL
	_this.keyTTL = tview.NewTextView()
	_this.keyTTL.SetBorder(true)
	_this.keyTTL.SetTitle("TTL")
	_this.keyTTL.SetTitleAlign(tview.AlignCenter)
	_this.keyInfoFlex.AddItem(_this.keyTTL, 3, 1, false)

	// 初始化KeyValue
	_this.KeyValue = tview.NewTextView()
	_this.KeyValue.SetBorder(true)
	_this.KeyValue.SetTitle("CONTENT")
	_this.KeyValue.SetTitleAlign(tview.AlignCenter)
	_this.KeyValue.SetFocusFunc(func() {
		// 设置当前焦点为键分组树
		_this.KeyValue.SetBorderColor(base.ActiveBorderColor)
	})
	_this.KeyValue.SetBlurFunc(func() {
		// 设置当前焦点为键分组树
		_this.KeyValue.SetBorderColor(base.InactiveBorderColor)
	})
	_this.keyInfoFlex.AddItem(_this.KeyValue, 0, 1, false)

	return nil
}

func (_this *RedisDataComponent) Start() {
	slog.Info("RedisDataComponent component start", "dbNum", _this.dbNum)
	_this.SetTreeData()
}

func (_this *RedisDataComponent) Stop() {

}

// --- data helpers ---

func _setTreeNodeData(node *tview.TreeNode, data *model.RedisGroupTree) {
	rootPath := node.GetReference()
	for _, child := range data.Children {
		fullPath := fmt.Sprintf("%s:%s", data.Name, child.Name)
		if rootPath == nil {
			// 如果有根路径，则拼接完整路径
			fullPath = child.Name
		} else {
			// 如果有根路径，则拼接完整路径
			fullPath = fmt.Sprintf("%s:%s", rootPath, child.Name)
		}

		childNode := tview.NewTreeNode(child.Name).
			SetColor(tcell.ColorGold).
			SetSelectable(true).
			SetExpanded(false)
		childNode.SetReference(fullPath)
		if len(child.Children) > 0 {
			childNode.SetColor(tcell.ColorGreen)
		}
		node.AddChild(childNode)
		_setTreeNodeData(childNode, child)
	}

}

// SetTreeData 设置树形数据
func (_this *RedisDataComponent) SetTreeData() {
	rootNode := tview.NewTreeNode("-").
		SetColor(tview.Styles.PrimaryTextColor)
	// 清空旧数据
	_this.keyGroupTree.SetRoot(rootNode)
	_this.keyGroupTree.SetCurrentNode(rootNode)
	_setTreeNodeData(rootNode, _this.redisData.Tree)

}

func NewRedisDataComponent(
	a *App,
	dbNum int,
	redisConnConfig *config.RedisConnConfig,
) *RedisDataComponent {
	var name = ""
	lp := RedisDataComponent{
		BaseFlex:        NewBaseFlex(name),
		app:             a,
		dbNum:           dbNum,
		redisConnConfig: redisConnConfig,
	}
	lp.SetDirection(tview.FlexRow)
	lp.SetBorder(false)

	return &lp
}
