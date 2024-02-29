#include "listuser.h"
#include "ui_listuser.h"

ListUser::ListUser(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::ListUser)
{
    ui->setupUi(this);
}

ListUser::~ListUser()
{
    delete ui;
}
