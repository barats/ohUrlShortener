$(document).ready(function() {

  $('.ui.dropdown').dropdown();

  $('#login-form')
    .form({
      fields: {
        account: {
          identifier  : 'account',
          rules: [
            {
              type   : 'empty',
              prompt : '账户名不能为空'
            },
            {
              type   : 'length[5]',
              prompt : '账户名长度不得少于5位'
            }
          ]
        },
        password: {
          identifier  : 'password',
          rules: [
            {
              type   : 'empty',
              prompt : '密码不能为空'
            },
            {
              type   : 'length[8]',
              prompt : '密码长度不得少于8位'
            }
          ]
        },
        captcha: {
          identifier  : 'captcha-text',
          rules: [
            {
              type   : 'empty',
              prompt : '验证码不能为空'
            },
            {
              type   : 'length[6]',
              prompt : '验证码长度不得少于6位'
            }
          ]
        }
      }
    });

    $('#form-search-url').form({
      fields:{
        url: {rules:[{
          type:'empty',
          prompt:'链接不能为空'
        }]}
      }
    });

    $('#logs-search-start-date').calendar({
      type:'date',
      text: {
        months: ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'],
        monthsShort: ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'],
      },
      formatter: {
        date: function (date, settings) {
          if (!date) return '';
          var day = date.getDate();
          var month = date.getMonth() + 1;
          var year = date.getFullYear();
          return year + '-' + month + '-' + day;
        }
      }
    });

  $('#logs-search-end-date').calendar({
    type:'date',
    text: {
      months: ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'],
      monthsShort: ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'],
    },
    formatter: {
      date: function (date, settings) {
        if (!date) return '';
        var day = date.getDate();
        var month = date.getMonth() + 1;
        var year = date.getFullYear();
        return year + '-' + month + '-' + day;
      }
    }
  });

    $('#sidebar-menu').sidebar('attach events', '#sidebar-menu-toggler');

    $('#btn-new-shorturl-modal').click(function(){
      $('#new-shorturl-modal').modal('show');
    });   
    
    $('#btn-gen-short-url').click(function() {
      var destUrl = $('#input_dest_url');
      var memo = $('#input_demo');
      var openType = $('#input_open_type')
      if( $.trim(destUrl.val()).length <= 0) {
        errorToast('目标链接不能为空！');
        destUrl.parent().addClass('error');
        return 
      }
      
      var data = {
        "dest_url": $.trim(destUrl.val()),
        "memo": $.trim(memo.val()),
        "open_type": $.trim(openType.val()),
      };

      $.ajax({
        type: "POST",
        url: '/admin/urls/generate',
        data: data,
        dataType: 'json',
        success: function() {                              
          successToast('新建成功！')
          destUrl.val('');
          memo.val('');
          $('#new-shorturl-modal').modal('hide'); 
        },
        error: function(e) {          
          errorToast($.parseJSON(e.responseText).message)
        } 
      });
    });//end of #btn-gen-short-url click
});

function successToast(message) {
  $('body').toast({
    class: 'success',
    displayTime: 2500,
    message: message,    
    showIcon:'exclamation circle',
    showProgress: 'bottom',
    onHidden: function() {location.reload()}
  });
}

function errorToast(message) {
  $('body').toast({
    class: 'error',
    displayTime: 2500,
    message: message,    
    showIcon:'exclamation circle',
    showProgress: 'bottom'
  });
}

function enable_url(url,state) {
  var data = {
    "dest_url": url,
    "enable": state
  }

  $.ajax({
    type: "POST",
    url: '/admin/urls/state',
    data: data,
    dataType: 'json',
    success: function() {                              
      successToast('操作成功')
    },
    error: function(e) {          
      errorToast($.parseJSON(e.responseText).message)
    } 
  });
}

function delete_url(url) {
  $('body').modal('confirm','温馨提示','确认删除该短链接吗？<br/>本操作将会同步删除该短链接的所有访问日志。', function(choice){
    if (choice) {
      $.ajax({
        type:"POST",
        url: "/admin/urls/delete",
        data : {
          "short_url": url,
        },
        success: function() {
          successToast('操作成功，再见！')
        },
        error: function(e) {
          errorToast($.parseJSON(e.responseText).message)
        }
      });
    }
  });
}

function copy_url(url){  
  navigator.clipboard.writeText(url);
  $('body').toast({
    class: 'success',
    displayTime: 2500,
    message: '复制成功',    
    showIcon:'exclamation circle',
    showProgress: 'bottom'
  });
}

function reload_captcha() {
  $.ajax({
    type: "POST",
    url: '/captcha',
    dataType: 'json',
    success: function(r) {            
      $('#captcha-image').html('<img src="/captcha/'+r.result+'.png" />');
      $('<input>').attr({type: 'hidden', value:r.result ,name: 'captcha-id'}).appendTo('#login-form');
    },
    error: function(e) {
      errorToast($.parseJSON(e.responseText).message)
    }
  });
}

function sign_out_config() {
  $('body').modal('confirm','温馨提示','确认退出 ohUrlShortener 管理后端吗？', function(choice){
    if (choice) {
      $.ajax({
        type:"POST",
        url: "/admin/logout",
        success: function() {
          successToast('操作成功，再见！')
        },
        error: function(e) {
          errorToast($.parseJSON(e.responseText).message)
        }
      });
    }
  });
}

function export_accesslog() {
  $('body').modal('confirm','温馨提示','确认导出访问日志?', function(choice){
    if (choice) {
      $("#form-export-logs").submit();
    }
  });
}