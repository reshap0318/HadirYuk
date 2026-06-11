package repositories

import "gorm.io/gorm"

type Repositories struct {
	TxManager             *TransactionManager
	User                  *UserRepository
	UserProfile           *UserProfileRepository
	PasswordReset         *PasswordResetRepository
	Permission            *PermissionRepository
	Role                  *RoleRepository
	RoleHasPerm           *RoleHasPermissionRepository
	UserRole              *UserRoleRepository
	Notification          *NotificationRepository
	Shift                 *ShiftRepository
	OfficeLocation        *OfficeLocationRepository
	LeaveType             *LeaveTypeRepository
	UserShiftAssignment   *UserShiftAssignmentRepository
	Attendance            *AttendanceRepository
	QRCode                *QRCodeRepository
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	txManager := NewTransactionManager(db)
	userRepo := NewUserRepository(db)
	userProfileRepo := NewUserProfileRepository(db)
	passwordResetRepo := NewPasswordResetRepository(db)
	permissionRepo := NewPermissionRepository(db)
	roleRepo := NewRoleRepository(db)
	roleHasPermRepo := NewRoleHasPermissionRepository(db)
	userRoleRepo := NewUserRoleRepository(db)
	notificationRepo := NewNotificationRepository(db)
	shiftRepo := NewShiftRepository(db)
	officeLocationRepo := NewOfficeLocationRepository(db)
	leaveTypeRepo := NewLeaveTypeRepository(db)
	userShiftAssignmentRepo := NewUserShiftAssignmentRepository(db)
	attendanceRepo := NewAttendanceRepository(db)
	qrCodeRepo := NewQRCodeRepository(db)

	return &Repositories{
		TxManager:             txManager,
		User:                  userRepo,
		UserProfile:           userProfileRepo,
		PasswordReset:         passwordResetRepo,
		Permission:            permissionRepo,
		Role:                  roleRepo,
		RoleHasPerm:           roleHasPermRepo,
		UserRole:              userRoleRepo,
		Notification:          notificationRepo,
		Shift:                 shiftRepo,
		OfficeLocation:        officeLocationRepo,
		LeaveType:             leaveTypeRepo,
		UserShiftAssignment:   userShiftAssignmentRepo,
		Attendance:            attendanceRepo,
		QRCode:                qrCodeRepo,
	}, nil
}
