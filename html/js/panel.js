/**
 * 递归查询，生成 json
 * @param {Number} master Master UID
 * @param {Object} data User Info List
 */
function recursive(master, data) {
  let masterTree = new Array
  masterTree.push({
    'name': data[master].username,
    'value': data[master].id,
    'children': []
  })
  for (let i in data) {
    const deep = function(tree) {
      for (let j = 0; j < tree.length; j++) {
        if (data[i].master == tree[j].value) {
          tree[j].children.push({
            'name': data[i].username,
            'value': data[i].id,
            'children': []
          })
          break
        }
        if (tree[j].children.length != 0) deep(tree[j].children)
      }
    }
    deep(masterTree)
  }
  return masterTree[0]
}

/**
 * 初始化树形图
 * @param {Object} obj
 */
function chartInit(obj, username) {
  let option = {
    tooltip: {
      trigger: 'item',
      triggerOn: 'mousemove'
    },
    series: [
      {
        type: 'tree',
        data: [obj],
        left: '2%',
        right: '2%',
        top: '8%',
        bottom: '20%',
        symbol: 'emptyCircle',
        orient: 'vertical',
        roam: true,
        expandAndCollapse: true,
        label: {
          normal: {
            position: 'insideBottom',
            offset: [0, -50],
            rotate: 0,
            verticalAlign: 'middle',
            align: 'right',
            fontSize: 15
          }
        },
        leaves: {
          label: {
            normal: {
              position: 'insideBottomRight',
              offset: [50, 0],
              rotate: -45,
              verticalAlign: 'middle',
              align: 'left',
              fontSize: 15
            }
          }
        },
        initialTreeDepth: -1,
        animationDurationUpdate: 750
      }
    ]
  }
  var myChart = echarts.init($('.app .chart'), null, { renderer: 'svg' })
  myChart.setOption(option, true)

  for (var i = 0; i < $('tspan').length; i++) {
    if ($('tspan')[i].innerHTML == username) {
      $('tspan')[i].setAttribute('class', 'target')
    } else {
      $('tspan')[i].removeAttribute('class')
    }
  }
}
