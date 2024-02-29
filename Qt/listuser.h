#ifndef LISTUSER_H
#define LISTUSER_H

#include <QWidget>

namespace Ui {
class ListUser;
}

class ListUser : public QWidget
{
    Q_OBJECT

public:
    explicit ListUser(QWidget *parent = nullptr);
    ~ListUser();

private:
    Ui::ListUser *ui;
};

#endif // LISTUSER_H
