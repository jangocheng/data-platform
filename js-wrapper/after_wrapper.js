input = {}
// json parse
input.Result = '{}'
input.Result = JSON.parse(input.Result)
console.log(input)

///////////////////////////////////////////////////////////////////////////////////////////////////

// company close filter
input = '{}'
targetArr = new Array(); newArr = new Array;  if (Array.isArray(input.Result)) { for ( var i in input.Result ) { if ( targetArr.includes(input.Result[i]["status"])) { newArr.push(input.Result[i]) } } };  input.Result = newArr;
console.log(input)

///////////////////////////////////////////////////////////////////////////////////////////////////

// default filter
input = '{}'
resultArr = new Array; resultArr.push({"Detail": input.Result}); input.Result =resultArr;
console.log(input)

// time_before filter
input = '{}'
resultArr = new Array;
if (Array.isArray(input.Result) && Array.isArray(input.Params.times)) {
    timeFlags = input.Params.times[input.Item.DeriveCode]; 
    if (timeFlags.length == 0) {
        timeFlags = new Array("-1");
    }
    for (var i in timeFlags) {
        var timeResult = new Array;
        var timeFlag = timeFlags[i];
        if (timeFlag == "-1") {
            resultArr.push({"Detail": input.Result, "Time": "-1"})
            continue
        }
        for ( var i in input.Result) {
            if (getPreTime(timeFlag)>input.Result[i]["PreDateCompare"]) {
                timeResult.push(input.Result[i])
            }
        }
        resultArr.push({"Detail": timeResult, "Time": timeFlag})
    }
    input.Result = resultArr;

} else {
    resultArr.push({"Detail": input.Result}); 
    input.Result =resultArr;
}
console.log(input)


// time_after filter
input = '{}'
resultArr = new Array;
if (Array.isArray(input.Result) && Array.isArray(input.Params.times)) {
    timeFlags = input.Params.times[input.Item.DeriveCode]; 
    if (timeFlags.length == 0) {
        timeFlags = new Array("-1");
    }
    for (var i in timeFlags) {
        var timeResult = new Array;
        var timeFlag = timeFlags[i];
        if (timeFlag == "-1") {
            resultArr.push({"Detail": input.Result, "Time": "-1"})
            continue
        }
        for ( var i in input.Result) {
            if (getPreTime(timeFlag)<input.Result[i]["LaterDateCompare"]) {
                timeResult.push(input.Result[i])
            }
        }
        resultArr.push({"Detail": timeResult, "Time": timeFlag})
    }
    input.Result = resultArr;

} else {
    resultArr.push({"Detail": input.Result}); 
    input.Result =resultArr;
}
console.log(input)


///////////////////////////////////////////////////////////////////////////////////////////////////

// 结果条数计算
input = '{}'

for (var i in input.Result) {
    input.Result[i]["Hit"] = 0 ? !input.Result[i].Detail : input.Result[i].Detail.length
}
console.log(input)

//////////////////////////////////////////////////////////////////////////////////////////////////

// 删除详情
detailFlag = input.Params.details || input.Params.details.includes(input.Item.DeriveCode); 
for (var i in input.Result) {
    if ( !detailFlag ) {
        delete input.Result[i]["Detail"] 
    }
}

