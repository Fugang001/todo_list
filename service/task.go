package service

import (
	"time"
	"todo_list/model"
	"todo_list/pkg/e"
	"todo_list/pkg/utils"
	"todo_list/serializer"
)

//创建任务的服务
type CreateTaskService struct {
	Title   string `form:"title" json:"title" binding:"required,min=2,max=100"`
	Content string `form:"content" json:"content" binding:"max=1000"`
	Status  int    `form:"status" json:"status"` //0 待办   1已完成
}

//展示任务详情的服务
type ShowTaskService struct {
}

//删除任务的服务
type DeleteTaskService struct {
}

//展示用户所有备忘录
type ListTasksService struct {
	Limit int `form:"limit" json:"limit"`
	Start int `form:"start" json:"start"`
}

//更新任务的服务
type UpdateTaskService struct {
	ID      uint   `form:"id" json:"id"`
	Title   string `form:"title" json:"title" binding:"required,min=2,max=100"`
	Content string `form:"content" json:"content" binding:"max=1000"`
	Status  int    `form:"status" json:"status"` //0 待办   1已完成
}

//搜索任务的服务
type SearchTaskService struct {
	Info     string `json:"info" form:"info"`
	PageNum  int    `json:"page_num" form:"page_num"`
	PageSize int    `json:"page_size" form:"page_size"`
}

//新增一条备忘录
func (service *CreateTaskService) Create(id uint) serializer.Response {
	var user model.User
	code := 200
	model.DB.First(&user, id) //查找user
	task := model.Task{       //把查到的数据赋值
		User:      user,
		Uid:       user.ID,
		Title:     service.Title,
		Content:   service.Content,
		Status:    service.Status,    //状态默认是0，未完成
		StartTime: time.Now().Unix(), //开始时间默认现在吧
		EndTime:   0,
	}
	//创建备忘录
	err := model.DB.Create(&task).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "创建备忘录失败",
			Error:  err.Error(),
		}
	}
	//没有错误
	return serializer.Response{
		Status: code,
		Msg:    "创建成功",
	}
}

//展示一条备忘录
func (service *ShowTaskService) Show(id string) serializer.Response {
	var task model.Task
	code := e.SUCCESS
	err := model.DB.First(&task, id).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildTask(task), //把任务序列化返回
		Msg:    e.GetMsg(code),
	}
}

//列表返回用户所有备忘录
func (service *ListTasksService) List(id uint) serializer.Response {
	var tasks []model.Task //是一个切片
	var total int64
	if service.Limit == 0 {
		service.Limit = 15 //默认15页
	}
	//多表查询	找到所有备忘录 		预加载User			找到是哪一个user			 计算出所有聚合函数
	model.DB.Model(model.Task{}).Preload("User").Where("uid = ?", id).Count(&total).
		Limit(service.Limit).Offset((service.Start - 1) * service.Limit).
		Find(&tasks)
	return serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(total))
}

//更新操作
func (service *UpdateTaskService) Update(id string) serializer.Response {
	var task model.Task
	model.DB.Model(model.Task{}).Where("id = ?", id).First(&task)
	task.Content = service.Content
	task.Status = service.Status
	task.Title = service.Title
	code := e.SUCCESS
	err := model.DB.Save(&task).Error
	if err != nil {
		utils.LogrusObj.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   "修改成功",
	}
}

//查询用户自己的备忘录操作，是一个模糊搜索
func (service *SearchTaskService) Search(uId uint) serializer.Response {
	var tasks []model.Task
	code := e.SUCCESS
	count := 0
	if service.PageSize == 0 {
		service.PageSize = 10
	}
	//model.DB.Where("uid=?", uId).Preload("User").First(&tasks) //先定位用户
	//再模糊查询
	err := model.DB.Model(&model.Task{}).Preload("User").Where("uid=?", uId).
		Where("title LIKE ? OR content LIKE ?", "%"+service.Info+"%", "%"+service.Info+"%").
		Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&tasks).Error
	if err != nil {
		utils.LogrusObj.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(count)),
	}
}

//删除备忘录
func (service *DeleteTaskService) Delete(id string) serializer.Response {
	var task model.Task
	code := e.SUCCESS
	err := model.DB.First(&task, id).Error //查找数据库找到备忘录
	if err != nil {
		utils.LogrusObj.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	err = model.DB.Delete(&task).Error //删除备忘录
	if err != nil {
		utils.LogrusObj.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    "删除成功。", //删除成功
	}
}
