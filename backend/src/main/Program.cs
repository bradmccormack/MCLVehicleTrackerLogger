using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Net.Sockets;
using System.Diagnostics;
using System.Xml;
using Mono.Data.SqliteClient;


namespace testwpf
{
    class Program
    {
        
	

	static void Main(string[] args)
        {
            UdpClient udpClient = new UdpClient("internal.myclublink.com.au", 6969);
             string connectionString = "URI=file:backend.db";

           IDbConnection dbcon;
           dbcon = (IDbConnection) new SqliteConnection(connectionString);
           dbcon.Open();
           IDbCommand dbcmd = dbcon.CreateCommand();

           string sql = "SELECT * FROM GPSRECORDS";
           dbcmd.CommandText = sql;
           IDataReader reader = dbcmd.ExecuteReader();
           while(reader.Read())
           {
                String data = reader.GetString(1) + ",L";
                data += reader.GetString(2) + ",";
                data += reader.GetString(3) + ",";
                data += "S" + reader.GetString(4) + ",";
                data += "H" + reader.GetString(5) + ",";
                data += "D" + reader.GetString(7) + ",";
                data += "F" + reader.GetString(6) == "1" ? "true" : "false" + ","
                data += reader.GetString(8);
                Byte[] sendBytes = Encoding.ASCII.GetBytes(data);
                /*
                Byte[] sendBytes = Encoding.ASCII.GetBytes("PHi there buddy,L" + lat +
                ",150.81094,S0.00,H147.2,D" + dt + ",Ftrue,BRADSBUS");
                */

                udpClient.Send(sendBytes, sendBytes.Length);
                System.Threading.Thread.Sleep(1000);
           }

        }
    }
}
