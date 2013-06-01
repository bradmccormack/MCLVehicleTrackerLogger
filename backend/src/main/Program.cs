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

            UdpClient udpClient = new UdpClient("127.0.0.1", 215);
            Byte[] sendBytes = Encoding.ASCII.GetBytes("PHi there buddy,L-34.50108,150.81094,S0.00,H147.2,D" + dt + ",Ftrue,ISUS01");
            try
            {
                udpClient.Send(sendBytes, sendBytes.Length);
                Console.ReadKey();

            }
            catch (Exception e)
            {
                Console.WriteLine(e.ToString());
            }
        }
    }
}
