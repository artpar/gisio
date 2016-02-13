/**
 * Created by fundo on 2/13/2016.
 */

function mapTo2dArray(data) {
    var keys = Object.keys(data);
    var arrayData = [];
    for (var i = 0; i < keys.length; i++) {
        arrayData.push([keys[i], data[keys[i]]]);
    }
    return arrayData;
}

function permutations(inputArr) {
    var results = [];

    function permute(arr, memo) {
        var cur, memo = memo || [];

        for (var i = 0; i < arr.length; i++) {
            cur = arr.splice(i, 1);
            if (arr.length === 0) {
                results.push(memo.concat(cur));
            }
            permute(arr.slice(), memo.concat(cur));
            arr.splice(i, 0, cur[0]);
        }

        return results;
    }

    return permute(inputArr);
}
/**
 * K-combinations
 *
 * Get k-sized combinations of elements in a set.
 *
 * Usage:
 *   k_combinations(set, k)
 *
 * Parameters:
 *   set: Array of objects of any type. They are treated as unique.
 *   k: size of combinations to search for.
 *
 * Return:
 *   Array of found combinations, size of a combination is k.
 *
 * Examples:
 *
 *   k_combinations([1, 2, 3], 1)
 *   -> [[1], [2], [3]]
 *
 *   k_combinations([1, 2, 3], 2)
 *   -> [[1,2], [1,3], [2, 3]
 *
 *   k_combinations([1, 2, 3], 3)
 *   -> [[1, 2, 3]]
 *
 *   k_combinations([1, 2, 3], 4)
 *   -> []
 *
 *   k_combinations([1, 2, 3], 0)
 *   -> []
 *
 *   k_combinations([1, 2, 3], -1)
 *   -> []
 *
 *   k_combinations([], 0)
 *   -> []
 */
function k_combinations(set, k) {
    var i, j, combs, head, tailcombs;

    if (k > set.length || k <= 0) {
        return [];
    }

    if (k == set.length) {
        return [set];
    }

    if (k == 1) {
        combs = [];
        for (i = 0; i < set.length; i++) {
            combs.push([set[i]]);
        }
        return combs;
    }

    // Assert {1 < k < set.length}

    combs = [];
    for (i = 0; i < set.length - k + 1; i++) {
        head = set.slice(i, i + 1);
        tailcombs = k_combinations(set.slice(i + 1), k - 1);
        for (j = 0; j < tailcombs.length; j++) {
            combs.push(head.concat(tailcombs[j]));
        }
    }
    return combs;
}


/**
 * Combinations
 *
 * Get all possible combinations of elements in a set.
 *
 * Usage:
 *   combinations(set)
 *
 * Examples:
 *
 *   combinations([1, 2, 3])
 *   -> [[1],[2],[3],[1,2],[1,3],[2,3],[1,2,3]]
 *
 *   combinations([1])
 *   -> [[1]]
 */
function combinations(set) {
    var k, i, combs, k_combs;
    combs = [];

    // Calculate all non-empty k-combinations
    for (k = 1; k <= set.length; k++) {
        k_combs = k_combinations(set, k);
        for (i = 0; i < k_combs.length; i++) {
            combs.push(k_combs[i]);
        }
    }
    return combs;
}