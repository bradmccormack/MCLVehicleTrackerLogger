using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Net.Sockets;
using System.Diagnostics;
using System.Xml;


namespace testwpf
{
    class Program
    {
        
	

	static void Main(string[] args)
        {
	    var dt = XmlConvert.ToString(DateTime.Now, XmlDateTimeSerializationMode.Utc);
            Console.WriteLine(dt);

            UdpClient udpClient = new UdpClient("127.0.0.1", 6969);	
            try
            {
		double lat = -34.50108;

		for(int i = 0; i < 10000 ; i ++)
		{
       		  Byte[] sendBytes = Encoding.ASCII.GetBytes("PHi there buddy,L" + lat + ",150.81094,S0.00,H147.2,D" + dt + ",Ftrue,ISUS01");
                  udpClient.Send(sendBytes, sendBytes.Length);
                  System.Threading.Thread.Sleep(1000);
		  lat += 0.00010f;
		}

            }
            catch (Exception e)
            {
                Console.WriteLine(e.ToString());
            }
        }
    }
}
