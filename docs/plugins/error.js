import { DOM, Fmt, style } from "./html.js";
import { parseCodeText, detectFileType, parseError } from "./parser.js";
import { search } from "./search.js";
import { SideBySide } from "./sidebyside.js";

let div = DOM("div");
let textarea = DOM("textarea");
let link = DOM("a");
let iframe = DOM("iframe");
let input = DOM("input");

let sbs = new SideBySide();

window.onload = () => {
    sbs.init();
};

async function popup() {
    let cont = div`class="${css.window}"`;
    // ${div`innerText="Run Code" class="${css.heading}"`}
    // ${textarea`class="${css.textarea}"`}
    let all = await getAll();
    let w = Fmt`${cont}
                    ${div`class="${css.bar}"`}
                        ${input`class="${css.searchBar}" placeholder="Search" type="search"`}
                    ${div`class="${css.debug}"`}
                        ${div`class="${css.inputs}"`}
                            ${div`innerText="Debug Error" class="${css.heading}"`}
                            ${textarea`class="${css.textarea}"`}

                        ${div`class="${css.submit}"`}
                            ${div`class="${css.button}" innerText="Submit"`.on(
                                "click",
                                (e) => {
                                    let ta = document.querySelectorAll(
                                        "." + css.textarea
                                    );
                                    if (ta[0].value != "") {
                                        // Debug
                                        debug(ta[0].value, cont, all);
                                    } else if (ta[1].value != "") {
                                        // Run
                                    }
                                }
                            )}}`;
    document.body.appendChild(w);
    document.querySelector(`.${css.searchBar}`).focus();
}

async function debug(code, cont, all) {
    cont.classList.toggle(css.fullscreen);
    cont.clear();
    let error = getError(code);

    let results = await search(detectFileType(code) + ": " + error);

    cont.add(searchResults(results));

    let errs = parseError(code);

    let func = findFunction(errs, all);

    let prevIndexs = [];
    prevIndexs.push(func);

    let frame = iframe`src="/${func.path}/?iframe=true&index=${
        func.index
    }#${func.name.toLowerCase()}${func.type}" class="${css.full}"`;

    frame.onload = () => {
        frame.contentWindow.postMessage(`ERROR: ${code}`);
    };
    cont.add(frame);

    window.onmessage = (e) => {
        if (e.data == "next") {
            let nextFunc = findFunction(
                errs,
                all,
                Math.min(func.index + 1, errs.length - 1)
            );

            if (func.index != nextFunc.index) {
                prevIndexs.push(func);
            }

            if (nextFunc.index != func.index) {
                func = nextFunc;
                frame.src = `/${func.path}/?iframe=true&index=${
                    func.index
                }#${func.name.toLowerCase()}${func.type}`;
            }
        } else if (e.data == "previous") {
            if (prevIndexs.length == 1) {
                func = prevIndexs[0];
            } else {
                func = prevIndexs.pop();
            }
            frame.src = `/${func.path}/?iframe=true&index=${
                func.index
            }#${func.name.toLowerCase()}${func.type}`;
        }
    };
}

function findFunction(errors, all, last = 0) {
    let base = findBasePath(errors);
    errors = errors.slice(last);
    let func = {};

    // console.log(errors, last, all);

    for (let a = 0; a < errors.length; a++) {
        const err = errors[a];
        let name = err.name.split(".").pop();
        let type = err.path.split(".").pop().trim();
        let path = err.path
            .slice(base.length)
            .split("/")
            .filter((e) => e.indexOf(".") == -1)
            .join("/");
        let bestKey = "";
        for (const key in all) {
            if (similarityScore(path, bestKey) < similarityScore(path, key)) {
                bestKey = key;
            }
        }
        // console.log(all[bestKey], bestKey, name, path, err.path);
        if (all[bestKey]) {
            let functions = all[bestKey].functions[type];
            if (functions) {
                let bestIndex = 0;
                for (let b = 0; b < functions.length; b++) {
                    if (
                        similarityScore(name, functions[bestIndex].name) <
                        similarityScore(name, functions[b].name)
                    ) {
                        bestIndex = b;
                    }
                }
                func = functions[bestIndex];
                func.type = type;
                func.path = bestKey;
                func.index = a + last;
                break;
            }
        }
    }
    return func;
}

function searchResults(results) {
    let d = div`class="${css.search}"`;
    for (let i = 0; i < results.length; i++) {
        const result = results[i];
        d.add(Fmt`${div`class="${css.row}"`}
                            ${div`class="${
                                css.title
                            }" innerText="${result.title.replace(/"/g, "")}"`}
                            ${div`class="${
                                css.description
                            }" innerText="${result.description.replace(
                                /"/g,
                                ""
                            )}"`}
                            ${div`class="${css.link}"`}
                                ${link`innerText="${result.url}" href="${result.url}" target="_blank"`}
                            `);
    }
    return d;
}

function getError(code) {
    let c = code.split("\n");
    return c[0];
}

async function getAll() {
    let homepage = await fetch("/").then((e) => e.text());
    let html = document.createElement("div");
    html.innerHTML = homepage;

    let pages = html.querySelectorAll("li a");

    let documents = {};

    for (let i = 0; i < pages.length; i++) {
        let url = new URL(pages[i].getAttribute("href"), window.location.origin)
            .href;
        let text = await fetch(url).then((e) => e.text());
        let frag = document.createElement("div");

        frag.innerHTML = text;
        let code = frag.querySelectorAll("pre code");

        let functions = {};

        for (let a = 0; a < code.length; a++) {
            let res = parseCodeText(code[a].innerText);
            if (functions[res.language] == undefined) {
                functions[res.language] = [];
            }
            functions[res.language].push(...res.functions);
        }

        documents[pages[i].getAttribute("href")] = { functions, code, frag };
    }

    console.log(documents);

    return documents;
}

function findBasePath(urls) {
    if (urls.length === 0) return ""; // If the array is empty, return an empty string

    // Split each path into an array of directories
    const paths = urls.map((url) => url.path.split("/"));

    // Find the shortest path length
    const minLength = Math.min(...paths.map((path) => path.length));

    let basePath = ""; // Initialize the base path

    // Iterate over each directory position
    for (let i = 0; i < minLength; i++) {
        const directory = paths[0][i]; // Get the directory at position i

        // Check if all paths have the same directory at position i
        if (paths.every((path) => path[i] === directory)) {
            basePath += directory + "/"; // Append the directory to the base path
        } else {
            break; // Stop iterating if there's a mismatch
        }
    }

    return basePath;
}

function similarityScore(string1, string2) {
    const matrix = Array.from({ length: string1.length + 1 }, (_, i) =>
        Array.from({ length: string2.length + 1 }, (_, j) =>
            i === 0 ? j : j === 0 ? i : 0
        )
    );

    for (let i = 1; i <= string1.length; i++) {
        for (let j = 1; j <= string2.length; j++) {
            if (string1[i - 1] === string2[j - 1]) {
                matrix[i][j] = matrix[i - 1][j - 1];
            } else {
                matrix[i][j] =
                    Math.min(
                        matrix[i - 1][j - 1],
                        matrix[i - 1][j],
                        matrix[i][j - 1]
                    ) + 1;
            }
        }
    }

    const maxLen = Math.max(string1.length, string2.length);
    const similarity = 1 - matrix[string1.length][string2.length] / maxLen;
    return similarity;
}

let css = style(/*css*/ `
    .window {
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        z-index: 100;
        width: 90%;
        max-width: 1100px;
        height: 90%;
        max-height: 800px;
    }
    .bar {
        width: 100%;
        position: absolute;
        top: 25%;
    }
    .searchBar {
        width: 100%;
        background: #26272b;
        border: none;
        height: 60px;
        padding: 10px;
        border-radius: 10px;
        color: #fff;
        font-size: 30px;
        box-shadow: 0 0 0 100000px rgba(0,0,0,.2);
    }
    .debug {
        background: #1a1b20;
        border-radius: 10px;
        height: 50%;
        position: absolute;
        width: 100%;
        bottom: 0;
    }
    .heading {
        color: #e7e7e7;
        font-family: sans-serif;
        font-size: 30px;
        font-weight: 600;
        margin-top: 30px;
    }
    .inputs {
        width: 90%;
        margin: auto;
        margin-top: 70px;
    }
    .textarea {
        width: -webkit-fill-available;
        height: 180px;
        background: #26272b;
        border: none;
        border-radius: 10px;
        margin-top: 20px;
        color: #fff;
        padding: 10px;
        font-family: monospace;
    }
    .submit {
        position: absolute;
        bottom: 15px;
        width: 90%;
        display: flex;
        justify-content: flex-end;
        margin: auto;
    }
    .button {
        width: 130px;
        height: 55px;
        background: #3a4db7;
        border-radius: 10px;
        line-height: 55px;
        text-align: center;
        color: #fff;
        font-family: sans-serif;
        font-weight: 700;
        cursor: pointer;
    }
    .fullscreen {
        width: 100vw;
        height: 100vh;
        max-height: 100vh;
        max-width: 100vw;
        top: 0;
        left: 0;
        transform: none;
        border-radius: 0px;
    }
    .search {
        width: 400px;
        color: #bcbcbc;
        font-family: sans-serif;
        float: right;
        background: #26272b;
        padding: 16px;
        height: 100%;
        overflow-y: auto;
    }
    .title {
        font-weight: 700;
        color: white;
    }
    .description {
        height: 55px;
        overflow-y: hidden;
    }
    .link {
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        }
    .row {
        margin-bottom: 20px;
        background: #1a1c20;
        padding: 12px;
        border-radius: 10px;
    }
    .full {
        width: calc(100vw - 432px);
        height: 100vh;
        border: none;
        background: #fff;
    }
    .error {
        background: #1a1c20;
        font-family: monospace;
        color: white;
        overflow: auto;
        padding: 13px;
        position: fixed;
        top: 0;
        right: 0;
        width: calc(50% - 26px);
        z-index: 1000;
        display: flex;
        flex-direction: column;
        align-items: flex-start;
    }
    .error > span {
        white-space: nowrap;
    }
    .highlight {
        background: rgba(210, 220, 0, 0.14);
    }
    .btnCont {
        display: flex;
        width: calc(100% - 20px);
        justify-content: space-between;
        background: #26272b;
        padding: 10px;
        border-radius: 10px;
        margin-top: 10px;
    }
    .btnCont > div {
        height: 35px;
        line-height: 35px;
        user-select: none;
    }
    .search {}
`);

(() => {
    if (window.location.search.indexOf("iframe") != -1) {
        window.onmessage = (event) => {
            if (event.data.slice(0, 7) == "ERROR: ") {
                const urlParams = new URLSearchParams(window.location.search);

                let text = event.data.slice(7);

                let errs = parseError(text);
                let highlight = errs[urlParams.get("index")].line;

                sbs.error(errs[urlParams.get("index")].lineNum);

                let err = div`class="${css.error}"`;
                err.innerHTML = text
                    .split("\n")
                    .map((e) => {
                        if (e == highlight) {
                            return `<span class="${css.highlight}">${highlight}</span>`;
                        } else {
                            return `<span>${e}</span>`;
                        }
                    })
                    .join("");
                err.add(Fmt`${div`class="${css.btnCont}"`}
                                ${div`class="${css.button}" innerText="Previous"`.on(
                                    "click",
                                    () => {
                                        window.top.postMessage("previous", "*");
                                    }
                                )}
                                ${div`class="${css.button}" innerText="Next"`.on(
                                    "click",
                                    () => {
                                        window.top.postMessage("next", "*");
                                    }
                                )}
                `);
                sbs.sideCont.parentElement.insertBefore(err, sbs.sideCont);
                sbs.sideCont.style.top = `${
                    parseInt(getComputedStyle(err)["height"]) + 40
                }px`;
            }
        };
    } else {
        let open = false;
        document.addEventListener("keydown", function (event) {
            // Check if Ctrl key is pressed on Windows or Command key on Mac
            const isCtrlOrCmdPressed =
                (event.ctrlKey && navigator.platform.indexOf("Win") > -1) ||
                (event.metaKey && navigator.platform.indexOf("Mac") > -1);

            // Check if 'K' key is pressed
            const isKPressed = event.key === "k" || event.keyCode === 75;

            // If both conditions are true, execute your code here
            if (isCtrlOrCmdPressed && isKPressed && !open) {
                event.preventDefault();
                open = true;
                // Your code here
                console.log("Ctrl or Command + K pressed");
                popup();
            } else if (event.key == "Escape") {
                let w = document.querySelector(`.${css.window}`);
                open = false;
                w.parentElement.removeChild(w);
            }
        });
    }
})();
