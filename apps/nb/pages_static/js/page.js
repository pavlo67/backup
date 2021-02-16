function object(A) {
    return A instanceof Object ? A : A === undefined || A === null ? {} : {'': A};
}

function array(A) {
    return A instanceof Array ? A : A === undefined || A === null ? [] : [A];
}

// function getData(formId) {
//     let data = {};
//     for (let el of Array.from(document.querySelectorAll(`[id^="${formId}"]`))) {
//         if (el.type === "checkbox" && !el.checked) {
//             continue;
//         }
//         data[el.id.substr(formId.length)] = el.value;
//     }
// }

// $(function() {
//     $("textarea[name='edit_text']").froalaEditor()
// });

// function message(text, isSuccess, clearContent) {
//     let elSuccess = document.getElementById('success_mes');
//     let elError = document.getElementById('error_mes');
//     if (!(elSuccess && elError)) return;
//
//     text = array(text);
//     let label = text.shift();
//     let html = label + "\n<p>" + text.map(t => "<li>" + t + "</li>\n").join("");
//
//     if (isSuccess) {
//         elSuccess.className = "alert success";
//         elSuccess.innerHTML = html;
//         elError.className = elError.innerText = "";
//     } else {
//         elError.className = "alert";
//         elError.innerHTML = html;
//         elSuccess.className = elSuccess.innerText = "";
//     }
//
//     if (clearContent) {
//         let el1 = document.getElementById('corpus');
//         if (el1) {
//             el1.innerText = '';
//         }
//     }
// }
//
// function saveMessageToLocalStorage(text, isSuccess) {
//     if (isSuccess) {
//         localStorage.setItem('success_message', text);
//     } else {
//         localStorage.setItem('error_message', text);
//     }
// }
//
// function showMessageFromLocalStorage() {
//     let elError = document.getElementById('error_mes');
//     let elSuccess = document.getElementById('success_mes');
//     let err = localStorage.getItem("error_message");
//     let success = localStorage.getItem("success_message");
//     if (err != null && err !== '') {
//         elError.className = "alert";
//         elError.innerText = err;
//         elSuccess.className = elSuccess.innerText = "";
//         localStorage.setItem('error_message', '');
//     } else if (success != null && success !== '') {
//         elSuccess.className = "alert success";
//         elSuccess.innerText = success;
//         elError.className = elError.innerText = "";
//         localStorage.setItem('success_message', '');
//     }
// }
//
// function tabSelect (id, groupID, idTitle, groupIDTitle) {
//     if (node = document.getElementById(id)) {
//         node.style.visibility = "visible";
//         node.style.position   = "static";
//         if (groupID !== "") {
//             elements = document.querySelectorAll(`[id^="${groupID}"]`);
//             for (let el of Array.from(elements)) {
//                 if (el.id !== id) {
//                     el.style.visibility = "hidden";
//                     el.style.position   = "absolute";
//                 }
//             }
//         }
//     }
//
//     if (nodeTitle = document.getElementById(idTitle)) {
//         nodeTitle.style.fontWeight = "bold";
//         if (groupIDTitle !== "") {
//             elements = document.querySelectorAll(`[id^="${groupIDTitle}"]`);
//             for (let el of Array.from(elements)) {
//                 if (el.id !== idTitle) {
//                     el.style.fontWeight = "normal";
//                 }
//             }
//         }
//     }
// }
//
// function openClose(imageID, contentID) {
//     let nodeImage = document.getElementById(imageID);
//     let nodeContent = document.getElementById(contentID);
//
//     if (nodeImage && nodeContent) {
//         if (nodeContent.style.visibility === "hidden") {
//             nodeContent.style.visibility = "inherit";
//             nodeContent.style.position = "static";
//             nodeImage.src = "/images/minus.png";
//         } else {
//             nodeContent.style.visibility = "hidden";
//             nodeContent.style.position = "absolute";
//             nodeImage.src = "/images/plus.png";
//         }
//     }
//
//
// }
