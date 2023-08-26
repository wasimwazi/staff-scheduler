package service

type Services struct {
	AccountService AccountAuthService
	AdminService   AdminService
	StaffService   StaffService
}

func NewServices(accountService AccountAuthService, adminService AdminService, staffService StaffService) *Services {
	return &Services{
		AccountService: accountService,
		AdminService:   adminService,
		StaffService:   staffService,
	}
}
