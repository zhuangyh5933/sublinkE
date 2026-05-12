package models

type SubLogs struct {
	ID            int
	IP            string
	Date          string
	Addr          string
	Region        string
	Client        string
	Status        string
	Count         int
	SubcriptionID int
	UserID        int
	Username      string
}

// Add 添加IP
func (iplog *SubLogs) Add() error {
	return DB.Create(iplog).Error
}

// 查找IP
func (iplog *SubLogs) Find(id int) error {
	query := DB.Where("ip = ? and subcription_id  = ?", iplog.IP, id)
	if iplog.UserID != 0 {
		query = query.Where("user_id = ?", iplog.UserID)
	}
	return query.First(iplog).Error
}

// Update 更新IP
func (iplog *SubLogs) Update() error {
	return DB.Where("id = ? or ip = ?", iplog.ID, iplog.IP).Updates(iplog).Error
}

// List 获取IP列表
func (iplog *SubLogs) List() ([]SubLogs, error) {
	var iplogs []SubLogs
	err := DB.Order("date desc").Find(&iplogs).Error
	if err != nil {
		return nil, err
	}
	return iplogs, nil
}

func (iplog *SubLogs) ListByUser(userID int) ([]SubLogs, error) {
	var iplogs []SubLogs
	err := DB.Where("user_id = ?", userID).Order("date desc").Find(&iplogs).Error
	if err != nil {
		return nil, err
	}
	return iplogs, nil
}

// Del 删除IP
func (iplog *SubLogs) Del() error {
	return DB.Delete(iplog).Error
}
