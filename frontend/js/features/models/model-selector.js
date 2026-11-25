import { getModels } from "../../api/api-deepseek.js";
import { DEFAULT_MODEL } from "../../constants.js";
import { modelSelect } from "../../ui/elements.js";
import { showError } from "../../ui/error-banner.js";

export function getSelectedModel() {
	return modelSelect.value || DEFAULT_MODEL;
}

function modelLabel(id) {
	if (id === "deepseek-chat") return "DeepSeek Chat";
	if (id === "deepseek-reasoner") return "DeepSeek Reasoner（推理）";
	return id;
}

export async function loadModels() {
	let ids = [];
	try {
		const models = await getModels();
		ids = (models || []).map((m) => m.id).filter(Boolean);
	} catch (error) {
		showError("模型列表获取失败，已回退到默认模型：" + error.message);
	}
	if (!ids.length) ids = [DEFAULT_MODEL];
	modelSelect.innerHTML = "";
	ids.forEach((id) => modelSelect.add(new Option(modelLabel(id), id)));
}

export function setModelValue(model) {
	if (!model) return;
	if (![...modelSelect.options].some((o) => o.value === model)) {
		modelSelect.add(new Option(modelLabel(model), model));
	}
	modelSelect.value = model;
}
