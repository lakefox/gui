function customDiff(text1, text2) {
    // Split the texts into arrays of lines
    var lines1 = text1.split("\n");
    var lines2 = text2.split("\n");

    var diff = [];
    var minLength = Math.min(lines1.length, lines2.length);

    // Compare each line of the texts
    for (var i = 0; i < minLength; i++) {
        if (lines1[i] !== lines2[i]) {
            // If lines are different, push the difference to the result array
            diff.push({
                value: lines1[i], // Line from the first text
                added: false, // Indicates that this line is not present in the second text
                removed: true, // Indicates that this line is removed from the first text
            });
            diff.push({
                value: lines2[i], // Line from the second text
                added: true, // Indicates that this line is not present in the first text
                removed: false, // Indicates that this line is added to the second text
            });
        } else {
            // If lines are the same, push a single line to the result array
            diff.push({
                value: lines1[i], // Line from both texts (same)
                added: false, // Indicates that this line is not added
                removed: false, // Indicates that this line is not removed
            });
        }
    }

    // If one text has more lines than the other, add the remaining lines to the result
    for (var j = minLength; j < lines1.length; j++) {
        diff.push({
            value: lines1[j],
            added: false,
            removed: true,
        });
    }

    for (var k = minLength; k < lines2.length; k++) {
        diff.push({
            value: lines2[k],
            added: true,
            removed: false,
        });
    }

    return diff;
}
