/**
 * It fetches the HTML of the DuckDuckGo search page, parses it, and returns the results
 * @param query - The query to search for
 * @returns An array of objects with the following properties:
 *     title: The title of the result
 *     description: The description of the result
 *     url: The url of the result
 */
export async function webSearch(query) {
    const html = await fetch(
        `https://cors.lowsh.workers.dev/?https://lite.duckduckgo.com/lite/?q=${encodeURIComponent(
            query
        )}`
    );
    let text = await html.text();
    let doc = parseHTML(text);
    let sponsored = [
        ...doc.querySelectorAll("tr[class='result-sponsored']"),
    ].pop();
    let trs = [...doc.querySelectorAll("tr")];
    let rawRes = [...chunks(trs.slice(trs.indexOf(sponsored) + 1), 4)];

    let results = [];
    for (let i = 0; i < rawRes.length; i++) {
        const group = rawRes[i];
        if (group.length == 4) {
            results.push({
                title: group[1].querySelector("a").textContent,
                description: group[2].querySelector(
                    "td[class='result-snippet']"
                ).textContent,
                url:
                    "http://" +
                    group[3].querySelector("span[class='link-text']")
                        .textContent,
            });
        }
    }
    return results;
}

export function searchDocuments(all, query, amount = 3) {
    let querySplit = query.split(" ");
    let newQuery = [];
    for (const key in all) {
        const page = all[key].frag.innerText
            .replace(/\n/g, " ")
            .split(" ")
            .filter((e) => e != "");
        for (let a = 0; a < page.length; a++) {
            const word = page[a];
            for (let b = 0; b < querySplit.length; b++) {
                if (newQuery[b] != undefined) {
                    if (
                        similarityScore(querySplit[b], word) >
                        similarityScore(querySplit[b], newQuery[b])
                    ) {
                        newQuery[b] = word;
                    }
                } else {
                    newQuery[b] = word;
                }
            }
        }
    }
    let topSentences = [];
    for (const key in all) {
        const sentences = all[key].frag.innerText
            .replace(/\n/g, ".")
            .split(".")
            .map((e) => e.trim())
            .filter((e) => e != "")
            .sort((a, b) => {
                let aScore = 0;
                let bScore = 0;
                for (let i = 0; i < newQuery.length; i++) {
                    if (a.indexOf(newQuery[i]) != -1) {
                        aScore += 1;
                    }
                    if (b.indexOf(newQuery[i]) != -1) {
                        bScore += 1;
                    }
                }
                return bScore - aScore;
            });
        for (let i = 0; i < amount; i++) {
            topSentences.push({ top: sentences[i], document: key });
        }
    }
    let top = topSentences.sort((a, b) => {
        let aScore = 0;
        let bScore = 0;
        for (let i = 0; i < newQuery.length; i++) {
            if (a.top.indexOf(newQuery[i]) != -1) {
                aScore += 1;
            }
            if (b.top.indexOf(newQuery[i]) != -1) {
                bScore += 1;
            }
        }
        return bScore - aScore;
    });
    let topDouble = top.slice(0, amount * 2);
    topDouble = topDouble.sort((a, b) => b.top.length - a.top.length);
    return topDouble.slice(0, amount);
}

export function similarityScore(string1, string2) {
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

function* chunks(arr, n) {
    for (let i = 0; i < arr.length; i += n) {
        yield arr.slice(i, i + n);
    }
}

function parseHTML(html) {
    const template = document.createElement("template");
    template.innerHTML = html;
    return template.content;
}
