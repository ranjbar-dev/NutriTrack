package shared

import "net/http"

// AppError is the centralized Persian error catalog.
// All API error responses MUST use entries from this catalog.
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"` // Always in Persian
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) ToResponse() map[string]any {
	return map[string]any{
		"code":    e.Code,
		"message": e.Message,
	}
}

// --- Internal errors ---
var ErrInternal = &AppError{
	Code:       "INTERNAL_ERROR",
	Message:    "خطای داخلی سرور رخ داده است",
	HTTPStatus: http.StatusInternalServerError,
}

var ErrNotFound = &AppError{
	Code:       "NOT_FOUND",
	Message:    "منبع مورد نظر یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

var ErrValidation = &AppError{
	Code:       "VALIDATION_ERROR",
	Message:    "اطلاعات ورودی نامعتبر است",
	HTTPStatus: http.StatusUnprocessableEntity,
}

var ErrUnauthorized = &AppError{
	Code:       "UNAUTHORIZED",
	Message:    "احراز هویت لازم است",
	HTTPStatus: http.StatusUnauthorized,
}

var ErrForbidden = &AppError{
	Code:       "FORBIDDEN",
	Message:    "شما دسترسی به این منبع را ندارید",
	HTTPStatus: http.StatusForbidden,
}

var ErrConflict = &AppError{
	Code:       "CONFLICT",
	Message:    "این اطلاعات قبلاً ثبت شده است",
	HTTPStatus: http.StatusConflict,
}

// --- Auth errors ---
var ErrInvalidCredentials = &AppError{
	Code:       "INVALID_CREDENTIALS",
	Message:    "ایمیل یا رمز عبور نادرست است",
	HTTPStatus: http.StatusUnauthorized,
}

var ErrInvalidToken = &AppError{
	Code:       "INVALID_TOKEN",
	Message:    "توکن نامعتبر یا منقضی شده است",
	HTTPStatus: http.StatusUnauthorized,
}

var ErrTokenRevoked = &AppError{
	Code:       "TOKEN_REVOKED",
	Message:    "توکن باطل شده است",
	HTTPStatus: http.StatusUnauthorized,
}

var ErrOTPInvalid = &AppError{
	Code:       "OTP_INVALID",
	Message:    "کد تأیید نادرست است",
	HTTPStatus: http.StatusUnauthorized,
}

var ErrOTPExpired = &AppError{
	Code:       "OTP_EXPIRED",
	Message:    "کد تأیید منقضی شده است",
	HTTPStatus: http.StatusUnauthorized,
}

var ErrOTPRateLimit = &AppError{
	Code:       "OTP_RATE_LIMIT",
	Message:    "تعداد درخواست‌های کد تأیید بیش از حد مجاز است. لطفاً چند دقیقه صبر کنید",
	HTTPStatus: http.StatusTooManyRequests,
}

var ErrOTPMaxAttempts = &AppError{
	Code:       "OTP_MAX_ATTEMPTS",
	Message:    "تعداد تلاش‌های نادرست بیش از حد مجاز است. کد جدید درخواست کنید",
	HTTPStatus: http.StatusUnauthorized,
}

// --- User errors ---
var ErrUserNotFound = &AppError{
	Code:       "USER_NOT_FOUND",
	Message:    "کاربر یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

var ErrUserAlreadyExists = &AppError{
	Code:       "USER_ALREADY_EXISTS",
	Message:    "کاربر با این اطلاعات قبلاً ثبت شده است",
	HTTPStatus: http.StatusConflict,
}

var ErrInvalidMobile = &AppError{
	Code:       "INVALID_MOBILE",
	Message:    "شماره موبایل نامعتبر است. فرمت صحیح: ۰۹xxxxxxxxx",
	HTTPStatus: http.StatusUnprocessableEntity,
}

// --- Food errors ---
var ErrFoodNotFound = &AppError{
	Code:       "FOOD_NOT_FOUND",
	Message:    "ماده غذایی یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

// --- Medication errors ---
var ErrMedicationNotFound = &AppError{
	Code:       "MEDICATION_NOT_FOUND",
	Message:    "دارو یا مکمل یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

// --- Diet plan errors ---
var ErrPlanNotFound = &AppError{
	Code:       "DIET_PLAN_NOT_FOUND",
	Message:    "برنامه غذایی یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

var ErrPlanAlreadyActive = &AppError{
	Code:       "PLAN_ALREADY_ACTIVE",
	Message:    "این کاربر از قبل یک برنامه غذایی فعال دارد",
	HTTPStatus: http.StatusConflict,
}

// --- Tracking errors ---
var ErrTrackingNotFound = &AppError{
	Code:       "TRACKING_NOT_FOUND",
	Message:    "رکورد پیگیری یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

// --- File errors ---
var ErrFileTooLarge = &AppError{
	Code:       "FILE_TOO_LARGE",
	Message:    "حجم فایل بیش از حد مجاز است",
	HTTPStatus: http.StatusRequestEntityTooLarge,
}

var ErrInvalidFileType = &AppError{
	Code:       "INVALID_FILE_TYPE",
	Message:    "نوع فایل مجاز نیست",
	HTTPStatus: http.StatusUnprocessableEntity,
}

// --- Lab result errors ---
var ErrLabResultNotFound = &AppError{
	Code:       "LAB_RESULT_NOT_FOUND",
	Message:    "نتیجه آزمایش یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

// --- Message errors ---
var ErrMessageNotFound = &AppError{
	Code:       "MESSAGE_NOT_FOUND",
	Message:    "پیام یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

// --- Food request errors ---
var ErrFoodRequestNotFound = &AppError{
	Code:       "FOOD_REQUEST_NOT_FOUND",
	Message:    "درخواست غذای مورد نظر یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

var ErrFoodRequestAlreadyProcessed = &AppError{
	Code:       "FOOD_REQUEST_ALREADY_PROCESSED",
	Message:    "این درخواست قبلاً پردازش شده است",
	HTTPStatus: http.StatusConflict,
}

var ErrFoodRequestNotOwned = &AppError{
	Code:       "FOOD_REQUEST_NOT_OWNED",
	Message:    "این درخواست متعلق به متخصص تغذیه شما نیست",
	HTTPStatus: http.StatusForbidden,
}

// --- Notification errors ---
var ErrNotificationPreferenceNotFound = &AppError{
	Code:       "NOTIFICATION_PREFERENCE_NOT_FOUND",
	Message:    "تنظیمات اعلان یافت نشد",
	HTTPStatus: http.StatusNotFound,
}

// WithMessage returns a new AppError with a custom Persian message (for dynamic error context)
func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    msg,
		HTTPStatus: e.HTTPStatus,
	}
}
