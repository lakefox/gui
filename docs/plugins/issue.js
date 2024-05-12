import { DOM, Fmt, style } from "./html.js";

let div = DOM("div");
let blockquote = DOM("blockquote");
let paragraph = DOM("p");
let h2 = DOM("h2");

let css = style(/*css*/ `
        .flag {
            border: 2px solid;
            padding: 10px;
            margin-bottom: 40px;
            max-width: 90%;
        }
        .row {
            border-bottom: 1px solid;
        }
    
    `);

(async () => {
    if (window.location.pathname == "/") {
        let docs = await getAll();
        console.log(docs);

        let anchors = document.querySelectorAll("li a");
        for (let i = 0; i < anchors.length; i++) {
            const element = anchors[i];

            let doc = docs[element.getAttribute("href")];

            let flagDefs = findMatches(doc.innerText);

            for (let b = 0; b < flagDefs.length; b++) {
                let flags = getFlags(doc, flagDefs[b].fullMatch);
                inject(element, flagDefs[b].word, flags, true);
            }
        }
    } else {
        let flagDefs = findMatches(document.body.innerText);

        for (let b = 0; b < flagDefs.length; b++) {
            let flags = getFlags(document, flagDefs[b].fullMatch);
            inject(document.querySelector("h1"), flagDefs[b].word, flags);
        }
    }
})();

function findMatches(text) {
    const regex = /\/\/\s*!([A-Za-z]+)\b:/g;
    const matches = [];
    const foundWords = new Set(); // Keep track of found words
    let match;

    while ((match = regex.exec(text)) !== null) {
        const word = match[1];
        if (!foundWords.has(word)) {
            // Check if word is already found
            matches.push({
                fullMatch: match[0],
                word: word,
            });
            foundWords.add(word); // Add word to found set
        }
    }

    return matches;
}

function getFlags(doc, flag = "// !ISSUE:") {
    return [...doc.querySelectorAll("pre code")]
        .map((el) => {
            return el.innerText
                .split("\n")
                .filter((e) => e.trim() != "")
                .filter((e) => parseInt(e).toString() != e)
                .map((e, i, a) => {
                    let s = [];
                    e = e.trim();
                    if (e.indexOf(flag) != -1) {
                        s.push(
                            e
                                .trim()
                                .slice(e.indexOf(flag) + 10)
                                .trim()
                        );
                        for (let k = i + 1; k < a.length; k++) {
                            if (a[k].trim().indexOf("// +") != -1) {
                                s.push(
                                    a[k]
                                        .trim()
                                        .slice(a[k].trim().indexOf("// +") + 4)
                                        .trim()
                                );
                            } else {
                                s.push(a[k].trim());
                                break;
                            }
                        }
                    }
                    return {
                        comment: s.slice(0, -1).join("\n"),
                        line: s.slice(-1)[0],
                    };
                })
                .filter((e) => e.comment != "");
        })
        .flat();
}

function inject(el, name, flags, home = false) {
    console.log(flags);
    if (flags.length > 0) {
        let html = div`class="${css.flag}"`;
        html.add(h2`innerText="${name.toUpperCase()}"`);

        for (let i = 0; i < flags.length; i++) {
            const flag = flags[i];
            if (home) {
                html.add(Fmt`${div`class="${css.row}"`}
                            ${div`innerText="${flag.comment.replace(
                                /\"/g,
                                ""
                            )}"`}
                    `);
            } else {
                html.add(Fmt`${div`class="${css.row}"`}
                    ${blockquote`id="${generateUniqueString(flag.line)}"`}
                            ${paragraph`innerText="${flag.line}"`}
                    ${div`innerText="${flag.comment.replace(/\"/g, "")}"`}
            `);
            }
        }

        el.insertAdjacentElement("afterend", html);
    }
}

function generateUniqueString(inputString) {
    // Generate a hash code from the input string
    let hash = 0;
    for (let i = 0; i < inputString.length; i++) {
        hash = (hash << 5) - hash + inputString.charCodeAt(i);
        hash &= hash; // Convert to 32bit integer
    }
    const hashCode = Math.abs(hash).toString(36);

    // Take the first 5 characters of the hash code to ensure consistency
    const uniqueString = hashCode.slice(0, 5);

    return uniqueString;
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

        documents[pages[i].getAttribute("href")] = frag;
    }

    return documents;
}
