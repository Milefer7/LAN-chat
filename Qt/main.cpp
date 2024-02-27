
#include "operation.h"

#include <QApplication>


int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    Operation w;
    w.show();
    return a.exec();
}
