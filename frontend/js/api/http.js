// 封装发HTTP请求的函数
// 输入(请求路径, 预期的请求体)，输出()
export async function request(path, options = {}) {
	const response = await fetch(path, {
		headers: { "Content-Type": "application/json" },
		...options,
	});

	let payload = null;
	
	try {
		payload = await response.json();
	} catch {
		payload = null;
	}

	if (!response.ok || !payload || payload.success === false) {
		let message
		
		// 有返回 且 返回的错误字段不为空 且 后端有返回消息 
		if(payload && payload.error && payload.error.message){
			message = payload.error.message;
		} else {
			message = `请求失败（HTTP ${response.status}）`;
		}
		
		throw new Error(message);
		
	}

	// 到这就成功了

	return payload.data;
}





// 读取错误信息
// 输入response，返回message(string)
export async function readErrorResponse(response) {
	let message;
	
	try {
		const payload = await response.json();
		
		if (payload && payload.error && payload.error.message) {
			message =  payload.error.message;
		} else {
			message = `请求失败（HTTP ${response.status}）`;
		}
	} catch {
		message = `请求失败（HTTP ${response.status}）`;
	}

	return message;
}
