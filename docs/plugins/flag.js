import { DOM, Fmt, style } from "./html.js";

(() => {
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

    let flags = document
        .querySelectorAll("pre code")[1]
        .innerText.split("\n")
        .filter((e) => e.trim() != "")
        .filter((e) => parseInt(e).toString() != e)
        .map((e, i, a) => {
            let s = [];
            if (e.indexOf("// !FLAG:") != -1) {
                s.push(e.trim().slice(9).trim());
                for (let k = i + 1; k < a.length; k++) {
                    if (a[k].indexOf("// +") != -1) {
                        s.push(a[k].trim().slice(4).trim());
                    } else {
                        s.push(a[k].trim());
                        break;
                    }
                }
            }
            return { comment: s.slice(0, -1).join("\n"), line: s.slice(-1)[0] };
        })
        .filter((e) => e.comment != "");

    if (flags.length > 0) {
        let title = document.querySelector("h1");

        let html = div`class="${css.flag}"`;
        html.add(h2`innerText="Flags!"`);

        for (let i = 0; i < flags.length; i++) {
            const flag = flags[i];
            html.add(Fmt`${div`class="${css.row}"`}
                    ${blockquote``}
                        ${paragraph`innerText="${flag.line}"`}
                    ${div`innerText="${flag.comment.replace(/\"/g, "")}"`}
            `);
        }

        title.insertAdjacentElement("afterend", html);
    }
})();
