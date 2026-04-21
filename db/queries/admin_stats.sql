-- name: GetAdminStats :one
SELECT
    COUNT(*) FILTER (WHERE role = 'nutritionist') AS total_nutritionists,
    COUNT(*) FILTER (WHERE role = 'nutritionist' AND is_active = true) AS active_nutritionists,
    COUNT(*) FILTER (WHERE role = 'nutritionist' AND is_active = false) AS inactive_nutritionists,
    COUNT(*) FILTER (WHERE role = 'client') AS total_clients,
    (SELECT COUNT(*) FROM foods) AS total_foods,
    (SELECT COUNT(*) FROM diet_plans WHERE is_active = true) AS active_diet_plans
FROM users;
