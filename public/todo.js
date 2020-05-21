(function ($) {
    'use strict';
    $(function () {
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');

        $('.todo-list-add-btn').on("click", function (event) {
            event.preventDefault();

            var inputVal = $(this).prevAll('.todo-list-input').val();

            if (inputVal) {
                // $.post("/todos", {name:inputVal}, function(e){
                //     addItem({name:inputVal, completed:false})
                // })
                $.post("/todoH", { name: inputVal }, addItem); //서버가 응답값 addItem으로 가도록

                todoListInput.val("");
            }

        });

        var addItem = function (item) {
            if (item.completed) {
                todoListItem.append("<li class='completed'" + " id='" + item.id + "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked'/>" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            } else {
                todoListItem.append("<li" + " id='" + item.id + "'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            }
        };

        $.get('/todoH', function (items) {
            items.forEach(e => {
                addItem(e);
            });
        });

        todoListItem.on('change', '.checkbox', function () {
            var id = $(this).closest("li").attr('id');
            var check = true;
            if ($(this).attr('checked'))
                check = false;

            var $self = $(this);
            $.get("complete-todoH/" + id + "?complete=" + check, function (data) {
                // if (data.success == false) {

                // }

                if (check) {
                    $self.attr('checked', 'checked');
                } else {
                    $self.removeAttr('checked');
                }
                $self.closest("li").toggleClass('completed');
            })
        });

        todoListItem.on('click', '.remove', function () {
            var id = $(this).closest("li").attr('id');
            var $self = $(this);
            $.ajax({
                url: "todoH/" + id,
                type: "DELETE",
                success: function (data) {
                    if (data.success) {
                        $self.parent().remove();
                    }
                }
            })

        });

    });
})(jQuery);