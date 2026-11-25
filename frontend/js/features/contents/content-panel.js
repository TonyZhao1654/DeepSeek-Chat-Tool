import {
	createContent,
	deleteContent as deleteContentApi,
	updateContent,
} from "../../api/api-contents.js";
import { contentsElement } from "../../ui/elements.js";
import { showError } from "../../ui/error-banner.js";
import { appState } from "../../state/app-state.js";

let pendingContentSave = Promise.resolve();

export function flushContentSaves() {
	return pendingContentSave;
}

export function renderContents(contents) {
	contentsElement.innerHTML = "";
	(contents || []).forEach((c) => contentsElement.appendChild(createContentElement(c)));
	refreshContentHint();
}

function refreshContentHint() {
	const hint = contentsElement.querySelector(".empty-hint");
	const hasContents = contentsElement.querySelector(".content-block");
	if (hasContents) {
		if (hint) hint.remove();
		return;
	}
	if (!hint) {
		const el = document.createElement("div");
		el.className = "empty-hint";
		el.textContent = "还没有内容块，点击右上角「新增内容块」手动添加。";
		contentsElement.appendChild(el);
	}
}

function createContentElement(content) {
	const wrapper = document.createElement("div");
	wrapper.className = "content-block";
	wrapper.dataset.id = content.id;

	const head = document.createElement("div");
	head.className = "content-head";

	const label = document.createElement("label");
	const check = document.createElement("input");
	check.type = "checkbox";
	check.checked = !!content.selected;
	label.appendChild(check);
	label.appendChild(document.createTextNode(" 启用"));
	head.appendChild(label);

	const removeBtn = document.createElement("button");
	removeBtn.className = "content-remove";
	removeBtn.textContent = "删除";
	head.appendChild(removeBtn);

	const titleInput = document.createElement("input");
	titleInput.type = "text";
	titleInput.className = "content-title";
	titleInput.placeholder = "标题";
	titleInput.value = content.title === "未命名内容块" ? "" : (content.title || "");

	const contentArea = document.createElement("textarea");
	contentArea.className = "content-text";
	contentArea.placeholder = "内容...";
	contentArea.value = content.content || "";

	wrapper.append(head, titleInput, contentArea);

	const save = () => saveContent(wrapper);
	check.addEventListener("change", save);
	titleInput.addEventListener("change", save);
	contentArea.addEventListener("change", save);
	removeBtn.addEventListener("click", () => deleteContent(wrapper, content.id));

	return wrapper;
}

export async function addContent() {
	if (!appState.sessionId) {
		showError("会话尚未就绪，请稍后再试");
		return;
	}
	try {
		const content = await createContent(appState.sessionId, {
			title: "",
			content: "",
			selected: true,
		});
		const hint = contentsElement.querySelector(".empty-hint");
		if (hint) hint.remove();
		const el = createContentElement(content);
		contentsElement.appendChild(el);
		el.querySelector(".content-title").focus();
		refreshContentHint();
	} catch (error) {
		showError("新增内容块失败：" + error.message);
	}
}

function saveContent(wrapper) {
	if (!appState.sessionId) return pendingContentSave;
	const sessionId = appState.sessionId;
	const body = {
		title: wrapper.querySelector(".content-title").value,
		content: wrapper.querySelector(".content-text").value,
		selected: wrapper.querySelector("input[type='checkbox']").checked,
	};
	pendingContentSave = pendingContentSave
		.then(() => updateContent(sessionId, Number(wrapper.dataset.id), body))
		.catch((error) => {
			showError("保存内容块失败：" + error.message);
		});
	return pendingContentSave;
}

async function deleteContent(wrapper, contentId) {
	if (!appState.sessionId) return;
	try {
		await deleteContentApi(appState.sessionId, Number(contentId));
		wrapper.remove();
		refreshContentHint();
	} catch (error) {
		showError("删除内容块失败：" + error.message);
	}
}
