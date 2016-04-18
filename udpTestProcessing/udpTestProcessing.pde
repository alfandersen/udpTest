import hypermedia.net.*; //UDP library by Stephane Cousot

UDP udp;  // define the UDP object

String remoteIP;
int remotePort;
boolean remoteFound = false; // only send to remote if we know where to send it to

byte val[] = new byte[2]; // values to send
int timeBetweenSends = 10; // send every 10 ms
long lastSend = 0; // last time send

void setup() {
  // create window and setup a simple coordinate system
  size(500,500);
  background(100);
  stroke(0);
  strokeWeight(3);
  line(0 , height/2 , width , height/2);
  line(width/2 , 0 , width/2 , height);
  
  // create a new datagram connection on port 6000
  // and wait for incomming message
  udp = new UDP( this, 6000 );
  udp.listen( true );
}

void draw() {
  if(mousePressed){
    calculateXY();
  }
  
  if(remoteFound && millis() > lastSend + timeBetweenSends){
    sendData();
    lastSend = millis();
  }
}

// map mouse values between -100 and 100
void calculateXY() {
  if(mouseX >= 0 && mouseX <= width)
    val[0] = byte(((mouseX - width/2.)/width)*200);
  if(mouseY >= 0 && mouseY <= height)
    val[1] = byte(((mouseY - height/2.)/height)*200);
}

// reset to zero when mouse button release
void mouseReleased() {
  val[0] = 0;
  val[1] = 0;
}

// send values as byte array via udp to the remote ip and port
void sendData(){
  udp.send( val, remoteIP, remotePort );
}

// this port recieved data from ip with port
void receive( byte[] data, String ip, int port ) {	// <-- extended handler ... void receive( byte[] data ) is the default
  // we now know the remote ip and port to send to
  if(!remoteFound){
    remoteFound = true;
    remoteIP = ip;
    remotePort = port;
  }
  
  // parse the data and write it out in the console
  data = subset(data, 0, data.length);
  String message = new String( data );
  println( "receive: \""+message+"\" from "+ip+" on port "+port );
}