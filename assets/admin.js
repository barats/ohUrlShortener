$(document).ready(function() {
  $('#login-form')
    .form({
      fields: {
        account: {
          identifier  : 'account',
          rules: [
            {
              type   : 'empty',
              prompt : 'Please enter account'
            },
            {
              type   : 'length[5]',
              prompt : 'Your account must be at least 5 characters'
            }
          ]
        },
        password: {
          identifier  : 'password',
          rules: [
            {
              type   : 'empty',
              prompt : 'Please enter your password'
            },
            {
              type   : 'length[8]',
              prompt : 'Your password must be at least 8 characters'
            }
          ]
        }
      }
    });
  

    $('#sidebar-menu').sidebar('attach events', '#sidebar-menu-toggler');

    $('#new-shorturl-btn').click(function(){
      $('#new-shorturl-modal').modal('show');
    });

});