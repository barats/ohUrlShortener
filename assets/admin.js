$(document).ready(function() {
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

    $('#form-search-logs').form({
      fields: {
        url: {rules:[{
          type:'empty',
          prompt:'链接不能为空'
        }]}
      }
    });
  

    $('#sidebar-menu').sidebar('attach events', '#sidebar-menu-toggler');

    $('#btn-new-shorturl-modal').click(function(){
      $('#new-shorturl-modal').modal('show');
    });   
    
    $('#btn-gen-short-url').click(function() {
      var destUrl = $('#input_dest_url');
      var memo = $('#input_demo');
      if( $.trim(destUrl.val()).length <= 0) {
        errorToast('目标链接不能为空！');
        destUrl.parent().addClass('error');
        return 
      }
      
      var data = {
        "dest_url": $.trim(destUrl.val()),
        "memo": $.trim(memo.val())
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
    displayTime: 3000,
    message: message,    
    showIcon:'exclamation circle',
    showProgress: 'bottom',
    onHidden: function() {location.reload()}
  });
}

function errorToast(message) {
  $('body').toast({
    class: 'error',
    displayTime: 3000,
    message: message,    
    showIcon:'exclamation circle',
    showProgress: 'bottom'
  });
}