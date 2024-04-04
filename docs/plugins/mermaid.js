import { parseCodeText } from "./parser.js";
import { DOM, Fmt, style } from "./html.js";

let blockquote = DOM("blockquote");
let p = DOM("p");
let table = DOM("table");
let thead = DOM("thead");
let tr = DOM("tr");
let th = DOM("th");
let td = DOM("td");

(() => {
    let code = document.querySelectorAll("pre code");

    let functions = {};

    for (let i = 0; i < code.length; i++) {
        let res = parseCodeText(code[i].innerText);

        if (functions[res.language] == undefined) {
            functions[res.language] = [];
        }
        functions[res.language].push(...res.functions);
    }

    let hs = document.querySelectorAll("h1,h2,h3,h4,h5,h6");

    let tagMatch = /\w+\?\((\w+)\)/;

    for (let i = 0; i < hs.length; i++) {
        const element = hs[i];
        if (tagMatch.test(element.innerText)) {
            let match = tagMatch.exec(element.innerText);
            let name = match[0].slice(0, match[0].indexOf("?"));
            let lang = match[1];
            element.innerHTML = element.innerHTML.slice(
                0,
                element.innerHTML.indexOf("?")
            );

            let func = functions[lang].filter((e) => {
                return e.name == name;
            })[0];

            // Inject blockquote

            let bq = Fmt`${blockquote``}
                            ${p`innerText="${func.definition}"`}`;
            element.parentElement.insertBefore(bq, element);

            let body = document.createElement("tbody");

            for (let i = 0; i < func.parameters.length; i++) {
                const element = func.parameters[i];
                body.appendChild(Fmt`${tr``}
                                        ${td`innerHTML="${element.name}"`}
                                        ${td`innerHTML="${element.type}"`}
                                        `);
            }

            body.appendChild(Fmt`${tr``}
                                        ${td`innerHTML="return"`}
                                        ${td`innerHTML="${func.return}"`}
                                        `);

            let t = Fmt`${table``}
                            ${thead``}
                                ${tr``}
                                    ${th`innerHTML="Name"`}
                                    ${th`innerHTML="Type"`}
                            ${body}`;

            element.insertAdjacentElement("afterend", t);
        }
    }
})();
