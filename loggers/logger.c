#include<stdio.h>
#include<string.h>    //strlen etc
#include<sys/socket.h> //sockets
#include<arpa/inet.h> //ipv4 addr
#include <unistd.h> //sleep
#include <sys/sysinfo.h> //memory free
#include <signal.h> //signals
#include <stdlib.h> //standard shit
#include <pthread.h> //posix threading
#include <semaphore.h> //signalling mechanism locking primative
#include <phidget21.h> //phidget stuff

//TODO move this to the header file along with function prototypes - This was a quick hack job to determine my capability of prototyping
//The original was in Java and that was really bad - Also flakey old version of OpenJDK on that Arm SOC so this is a feasbility hack :-)

struct LoggerDeviceInfo {
    CPhidgetInterfaceKitHandle IfKit;
    CPhidgetGPSHandle GPS;

} LoggerDevice;


struct LoggerInfo {
    double LastLatitude;
    double LastLongitude;
    double Altitude;
    double Speed;
    double Heading;
    char DateTime[25];
    int Fix;
    GPSDate LastDate;
    GPSTime LastTime;
    uint isRunning;
    ulong memFree;
    char ID[255];
} Logger;

//semaphore for signalling other thread
static sem_t s_sem;

//mutex for locking shared LoggerDeviceInfo and LoggerInfo structs
static pthread_mutex_t lock, devicelock;


long getFreeRam()
{
    //parse /proc/meminfo so we can get memfree + cache amount so it gives real free memory value

    FILE* fp = fopen( "/proc/meminfo", "r" );
    if ( fp != NULL )
    {
        size_t bufsize = 1024 * sizeof(char);
        char* buf      = (char*)malloc( bufsize );
        long memfree     = -1L;
        long membuffered = -1L;

        while ( getline( &buf, &bufsize, fp ) >= 0 )
        {
            if( strncmp(buf, "MemFree", 7) == 0) {
                sscanf(buf, "%*s%ld", &memfree);
            }
            else if(strncmp(buf, "Cached", 6) == 0) {
                sscanf(buf, "%*s%ld", &membuffered);
            }
            else
                continue;
        }
        fclose( fp );
        free( (void*)buf );
        return (memfree + membuffered) * 1024;
    }

    return -1;
}


void sig_handler(int signo)
{
    if(signo == SIGINT)
    {
        pthread_mutex_lock(&lock);

        Logger.isRunning = 0;

        pthread_mutex_unlock(&lock);

        puts("\nExiting .. control c trapped\n");
    }
}


void *transmitter(void * ptr)
{

    int socket_desc;
    struct sockaddr_in server;

    struct LoggerInfo *loggerInfo = (struct LoggerInfo*) ptr;

    //Create socket
    socket_desc = socket(AF_INET , SOCK_STREAM , 0);
    if (socket_desc == -1)
    {
        puts("\nCould not create TCP socket\n");
    }

    //TODO remove the hardcoded address and pass the port and ip in via a struct and de-reference via ptr
    //this is an example - original removed
    server.sin_addr.s_addr = inet_addr("1.2.3.4");
    server.sin_family = AF_INET;
    server.sin_port = htons( 6969 );

    puts("\nConnecting to backend ...\n");

    //TODO this needs to be refactored .. if connection fails it doesn't retry.
    if (connect(socket_desc , (struct sockaddr *)&server , sizeof(server)) < 0)
    {
        puts("\nCannot connect to remote backend !\n");
        exit(1);
    }

    puts("\nConnected to backend\n");

    char *json = malloc(sizeof(struct LoggerInfo) + 100); // figure out the JSON encoding width and add it.. 100 for now
    if(json == 0)
    {
        puts("Failed to allocate memory for JSON packet\n");
        exit(-1);
    }

    while(loggerInfo->isRunning)
    {
        sem_wait(&s_sem);

        pthread_mutex_lock(&lock);

        //critical section stuff
        Logger.memFree = getFreeRam();
        pthread_mutex_unlock(&lock);

        int len = snprintf(json, sizeof(struct LoggerInfo) + 100, "{\"sentence\":\"PL%f,%f,S%f,H%f,D2014-10-23T02:52:59.132Z,Ftrue,%s\",\"diagnostics\":\"CT69.69,CV69.700001,CF400,MF%lu\"}",
            Logger.LastLatitude, Logger.LastLongitude, Logger.Speed, Logger.Heading, Logger.ID, Logger.memFree);

        if(len == 0)
        {
            puts("Failed to format the JSON packet using sprintfn\n");
            free(json);
            continue;
        }

        if( send(socket_desc , json , strlen(json) , 0) < 0)
            puts("Send Failed\n");

   }

    free(json);

    puts("\n Transmitter is joining main thread \n");
    puts("\n Exiting logger.. ");
    shutdown(socket_desc, 2);
    return (void*) 0;
}




int CCONV GPSErrorHandler(CPhidgetHandle phid, void *userptr, int ErrorCode, const char *unknown)
{
    printf("Error handler ran!\n");
    return 0;
}

int CCONV GPSposnChange(CPhidgetGPSHandle phid, void *userPtr, double latitude, double longitude, double altitude)
{
    GPSDate date;
    GPSTime time;
    CPhidgetGPSHandle gps = (CPhidgetGPSHandle)phid;

    pthread_mutex_lock(&lock);

    Logger.LastLatitude = latitude;
    Logger.LastLongitude = longitude;
    Logger.Altitude = altitude;


    printf("Position Change event: lat: %3.4lf, long: %4.4lf, alt: %5.4lf\n", latitude, longitude, altitude);

    if(!CPhidgetGPS_getDate(gps, &Logger.LastDate) && !CPhidgetGPS_getTime(gps, &Logger.LastTime))
        printf(" Date: %02d/%02d/%02d Time %02d:%02d:%02d.%03d\n", date.tm_mday, date.tm_mon, date.tm_year, time.tm_hour, time.tm_min, time.tm_sec, time.tm_ms);

    if(!CPhidgetGPS_getHeading(gps, &Logger.Heading) && !CPhidgetGPS_getVelocity(gps, &Logger.Speed))
        printf(" Heading: %3.2lf, Velocity: %4.3lf\n",Logger.Heading, Logger.Speed);


    pthread_mutex_unlock(&lock);
    return 0;
}

int CCONV GPSfixChange(CPhidgetGPSHandle phid, void *userPtr, int status)
{
    //Need to update the current fix value

    pthread_mutex_lock(&lock);

    //Transmitter might be touching this so lock
    Logger.Fix = status;

    pthread_mutex_unlock(&lock);

    return 0;
}


int CCONV AttachHandler(CPhidgetHandle phid, void *userPtr)
{
    CPhidget_DeviceClass cls;

    CPhidget_getDeviceClass(phid, &cls);

    //http://www.phidgets.com/documentation/web/cdoc/group__phidconst.html
    if(cls == PHIDCLASS_GPS)
    {
        puts("GPS Attached \n");
        pthread_mutex_lock(&devicelock);

        LoggerDevice.GPS = (CPhidgetGPSHandle) phid;

        CPhidgetGPS_create(&LoggerDevice.GPS);

        CPhidget_set_OnError_Handler((CPhidgetHandle) LoggerDevice.GPS, GPSErrorHandler, NULL);
        CPhidgetGPS_set_OnPositionChange_Handler(LoggerDevice.GPS, GPSposnChange, NULL);
        CPhidgetGPS_set_OnPositionFixStatusChange_Handler(LoggerDevice.GPS, GPSfixChange, NULL);

        puts("About to open device \n");
        CPhidget_open((CPhidgetHandle) LoggerDevice.GPS, -1);
        puts("Device Open complete\n");

        pthread_mutex_unlock(&devicelock);

    }
    if(cls == PHIDCLASS_INTERFACEKIT)
    {
        puts("Interface Kit Attached\n");
    }


    return 0;
}

int CCONV DetachHandler(CPhidgetHandle phid, void *userPtr)
{
    int serialNo;
    const char *name;
    CPhidget_DeviceClass cls;

    CPhidget_getDeviceName (phid, &name);
    CPhidget_getSerialNumber(phid, &serialNo);
    CPhidget_getDeviceClass(phid, &cls);
    printf("%s %10d detached!\n", name, serialNo);


    if(cls == PHIDCLASS_GPS)
    {
        puts("GPS Detached \n");
        pthread_mutex_lock(&devicelock);

        CPhidget_close((CPhidgetHandle) LoggerDevice.GPS);
        CPhidget_delete((CPhidgetHandle) LoggerDevice.GPS);
        LoggerDevice.GPS = NULL;

        pthread_mutex_unlock(&devicelock);

        //lock a mutex, modify the phidgetdevice struct and clear the member -- same for attach
    }
    if(cls == PHIDCLASS_INTERFACEKIT)
    {
        puts("Interface Kit Detached\n");

    }

    return 0;
}

int CCONV ErrorHandler(CPhidgetManagerHandle MAN, void *usrptr, int Code, const char *Description)
{
    printf("Error handled. %d - %s\n", Code, Description);
    return 0;
}


int main(int argc , char *argv[])
{
    CPhidgetManagerHandle man = 0;
    CPhidget_enableLogging(PHIDGET_LOG_VERBOSE, NULL);
    CPhidgetManager_create(&man);

    //Set the handlers to be run when the device is plugged in or opened from software, unplugged or closed from software, or generates an error.
    CPhidgetManager_set_OnAttach_Handler(man, AttachHandler, man);
    CPhidgetManager_set_OnDetach_Handler(man, DetachHandler, man);
    CPhidgetManager_set_OnError_Handler(man, ErrorHandler, NULL);

    //open the Manager for device connections
    CPhidgetManager_open(man);

    pthread_t thread;

    //init semaphore and mutex
    if(sem_init(&s_sem, 0, 0) != 0)
    {
        puts("\n Failed to initialise semphore\n");
        return -1;
    }
    if(pthread_mutex_init(&lock, NULL) != 0)
    {
        puts("\n Failed to initialize mutex\n");
        return -1;
    }

    //TODO remove this.. its just for testing
    Logger.LastLatitude=0.0f;
    Logger.LastLongitude = 0.0f;
    Logger.Speed = 0.0f;
    Logger.Heading = 0.0f;
    Logger.Fix = 0;
    Logger.isRunning = 1;

    // remove this hard coding
    strncpy(Logger.ID, "CLogger", strlen("Mr C Logger"));

    if(pthread_create(&thread, NULL, transmitter, &Logger) != 0)
    {
        puts("\n Failed to create tcp sending thread\n");
        return -1;
    }


    //catch control c
    if(signal(SIGINT, sig_handler) == SIG_ERR)
        puts("Cannot catch SIGINT\n");

    while(Logger.isRunning)
    {
        sleep(1);
        //TODO more logic will occur in main thread later ..

        //signal semaphore
        sem_post(&s_sem);
    }

    //wait for thread to join main thread
    puts("\n\nCleaning up ..\n");
    sem_post(&s_sem);
    pthread_join(thread, NULL);

    //clean up semaphore and mutex
    sem_destroy(&s_sem);
    pthread_mutex_destroy(&lock);

    //clean up Phidget stuff
    CPhidgetManager_close(man);
    CPhidgetManager_delete(man);

    return 0;
}
