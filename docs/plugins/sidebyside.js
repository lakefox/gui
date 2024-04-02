import { DOM, style, inline } from "./html.js";

let div = DOM("div");

export class SideBySide {
    constructor() {
        this.body = document.querySelector("body");
        this.sideCont = div`class="${css.container}"`;
        this.pre = document.querySelectorAll("pre");
        this.main = document.querySelector("main");
        this.lineLocked = false;
    }

    init() {
        this.main.classList.add(css.main);
        this.body.classList.add(css.body);
        this.body.appendChild(this.sideCont);
        let markedItems = [...document.querySelectorAll("blockquote > p")];
        this.markedItems = markedItems;
        let documents = [...this.pre].map(getText);
        this.highlights = [];
        this.current = 0;

        for (let i = 0; i < markedItems.length; i++) {
            const text = markedItems[i].innerText;
            let found = find(text, documents);
            if (found[0] == 0 && found[1] == Infinity && found[0] == 0) {
                markedItems.splice(i, 1);
                i--;
            } else {
                this.highlights.push(found);
                this.sideCont.appendChild(this.pre[found[2]]);
                inline(this.pre[found[2]], { display: "none" });
                markedItems[i].dataset.index = this.highlights.length - 1;
                inline(markedItems[i], {
                    color: "transparent",
                    height: "0px",
                    overflow: "hidden",
                });
            }
        }
        inline(this.pre[this.highlights[0][2]], { display: "block" });

        let currentItem = -1;
        this.main.onscroll = (e) => {
            if (!this.lineLocked) {
                let closestItem;
                let minDistance = Infinity;

                markedItems.forEach((item, i) => {
                    if (isInViewport(item) || i == markedItems.length - 1) {
                        const distance = Math.abs(
                            item.getBoundingClientRect().top
                        );
                        if (distance < minDistance) {
                            minDistance = distance;
                            closestItem = item;
                        }
                    }
                });

                let found =
                    this.highlights[parseInt(closestItem.dataset.index)];
                if (currentItem != parseInt(closestItem.dataset.index)) {
                    for (let i = 0; i < this.highlights.length; i++) {
                        const element = this.highlights[i];
                        unHighlight(this.pre[element[2]], element);
                        inline(this.pre[element[2]], { display: "none" });
                    }

                    inline(this.pre[found[2]], { display: "block" });
                    highlight(this.pre[found[2]], found);
                    view(this.sideCont, this.pre[found[2]], found);
                    currentItem = parseInt(closestItem.dataset.index);
                }
            }
        };

        this.main.appendChild(div`style="padding: 100vh;"`);
    }

    error(start, end = false) {
        if (end == false) {
            end = start;
        }

        // this.lineLocked = true;

        let closestItem;
        let minDistance = Infinity;

        this.markedItems.forEach((item, i) => {
            if (isInViewport(item) || i == this.markedItems.length - 1) {
                const distance = Math.abs(item.getBoundingClientRect().top);
                if (distance < minDistance) {
                    minDistance = distance;
                    closestItem = item;
                }
            }
        });
        let found = this.highlights[parseInt(closestItem.dataset.index)];
        found[0] = start;
        found[1] = end;
        highlight(this.pre[found[2]], found, "#754242");
        view(this.sideCont, this.pre[found[2]], found);
    }
}

let css = style(/*css*/ `
    .main {
        width: 50%;
        overflow-y: auto;
        height: calc(100vh - 156px);
        margin: 0;
        padding-bottom: 100px;
    }
    .body {
        overflow-y: hidden;
    }
    .container {
        width: 50%;
        position: fixed;
        right: 0;
        top: 0;
        height: 100vh;
        overflow-y: auto;
    }
`);

function getText(el) {
    let t = "";
    let code = el.querySelector("code");

    if (code) {
        for (let a = 0; a < code.children.length; a++) {
            const element = code.children[a];
            t += element.children[1].innerText;
        }
    }
    return t;
}

function find(text, documents) {
    let q = text.replace(/\s/g, "");

    let index = 0;
    let start = 0;
    let end = Infinity;

    for (let i = 0; i < documents.length; i++) {
        const d = documents[i].replace(/\s/g, "");
        if (d.indexOf(q) != -1) {
            let lines = documents[i].split("\n");
            for (let a = 0; a < lines.length; a++) {
                if (
                    lines
                        .slice(0, a)
                        .join("\n")
                        .replace(/\s/g, "")
                        .indexOf(q) != -1
                ) {
                    end = a;
                    for (let b = 0; b < a + 1; b++) {
                        if (
                            lines
                                .slice(b, a)
                                .join("\n")
                                .replace(/\s/g, "")
                                .indexOf(q) == -1
                        ) {
                            start = b;
                            index = i;
                            break;
                        }
                    }
                    break;
                }
            }
        }
    }
    return [start, end, index];
}

function highlight(el, lines, color = "#d2dc0024") {
    let code = el.querySelector("code");

    for (let a = lines[0] - 1; a < lines[1]; a++) {
        const element = code.children[a];
        inline(element, {
            background: color,
        });
    }
}

function unHighlight(el, lines) {
    let code = el.querySelector("code");

    for (let a = lines[0] - 1; a < lines[1]; a++) {
        const element = code.children[a];
        inline(element, {
            background: "",
        });
    }
}

function isInViewport(element) {
    const rect = element.getBoundingClientRect();
    return (
        rect.top >= 0 &&
        rect.left >= 0 &&
        rect.bottom <=
            (window.innerHeight || document.documentElement.clientHeight) &&
        rect.right <=
            (window.innerWidth || document.documentElement.clientWidth)
    );
}

function view(parentElement, el, lines) {
    let code = el.querySelector("code");
    let childElement = code.children[lines[0] - 1];
    const parentRect = parentElement.getBoundingClientRect();
    const childRect = childElement.getBoundingClientRect();
    const parentScrollTop = parentElement.scrollTop;

    // Calculate the vertical position to scroll to center the child element
    const targetScrollTop =
        childRect.top -
        parentRect.top -
        (parentRect.height - childRect.height) / 2;

    // Scroll the parent element to the calculated position
    parentElement.scrollTo({
        top: parentScrollTop + targetScrollTop,
        behavior: "smooth", // Optional: for smooth scrolling
    });
}
