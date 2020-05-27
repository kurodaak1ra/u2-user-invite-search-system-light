/**
 * 仿 jq 选择器
 * @param {String} ele 名称
 */
function $(ele) {
  let el = document.querySelectorAll(ele).length
  if (el === 0) {
    return null
  } else if (el === 1) {
    return document.querySelectorAll(ele)[0]
  } else {
    return document.querySelectorAll(ele)
  }
}

/**
 * 弹窗开
 * @param {String} el 类名
 */
function show(el) {
  if (el == '.menu') {
    $('.menu').css('left', '0')
    setTimeout(function() {
      $('.menu input[type="text"]').focus()
    })
  } else {
    hide('.menu')
    $(`${el}`).css('display', 'table')
  }
}

/**
 * 弹窗关
 * @param {String} el 类名
 */
function hide(el) {
  if (el == '.menu') {
    $('.menu').css('left', /Andriod|iPhone|iPad/i.test(navigator.userAgent) ? '-100%' : '-450px')
  } else {
    $(`${el}`).css('display', 'none')
  }
}

/**
 * Prototype
 */
Element.prototype.on = function(way, func) {
  this.addEventListener(way, func)
}
Element.prototype.css = function(key, val) {
  this.style[key] = val
}

