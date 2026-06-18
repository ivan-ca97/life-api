package permissions

// Permission constants follow the convention "{resource}:{action}".
const (
	UsersRead       = "users:read"
	UsersCreate     = "users:create"
	UsersUpdate     = "users:update"
	UsersDeactivate = "users:deactivate"

	MealsRead   = "meals:read"
	MealsCreate = "meals:create"
	MealsUpdate = "meals:update"
	MealsDelete = "meals:delete"

	ExercisesRead   = "exercises:read"
	ExercisesCreate = "exercises:create"
	ExercisesUpdate = "exercises:update"
	ExercisesDelete = "exercises:delete"

	FoodsRead   = "foods:read"
	FoodsCreate = "foods:create"
	FoodsUpdate = "foods:update"
	FoodsDelete = "foods:delete"

	DailyRead   = "daily:read"
	DailyUpdate = "daily:update"

	WeightRead   = "weight:read"
	WeightCreate = "weight:create"
	WeightUpdate = "weight:update"
	WeightDelete = "weight:delete"

	GoalsRead   = "goals:read"
	GoalsUpdate = "goals:update"

	SharesRead   = "shares:read"
	SharesCreate = "shares:create"
	SharesDelete = "shares:delete"

	StepsRead  = "steps:read"
	StepsWrite = "steps:write"
)
