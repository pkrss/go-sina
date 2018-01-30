#include <stdio.h>
#include "../quote.h"

// force gcc to link in go runtime (may be a better solution than this)
void dummy() {
    SoSetOnPublishChannel(0);
    SoSetOnPublishChannelC(0);
	GoString g;
    SoNotifyStkChanged(g);
    SoNotifyStkChangedC(0);
	SoStart();
    SoRunOnce();
}

int main() {

}