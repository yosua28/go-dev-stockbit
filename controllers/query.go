package controllers

var query string = `SELECT 
						u.id AS ID, 
						u.user_name AS UserName,
						up.user_name AS ParentUserName
					FROM user AS u
					LEFT JOIN user AS up ON up.id = u.parent`
