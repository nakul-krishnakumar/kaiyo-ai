import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import { MapPin, Hotel, Utensils, Plane } from "lucide-react";
import "leaflet/dist/leaflet.css";
import L from "leaflet";

// Fix for default markers in React Leaflet
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
  iconUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
  shadowUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
});

export interface ItineraryDay {
  day: number;
  title: string;
  activities: string[];
}

export interface TravelData {
  destination: string;
  totalCost: string;
  dates: string;
  coordinates: [number, number];
  itinerary: ItineraryDay[];
}

interface TravelResultsProps {
  travelData: TravelData | null;
}

export const TravelResults = ({ travelData }: TravelResultsProps) => {
  if (!travelData) {
    return (
      <div className="w-96 bg-card border-l border-border p-6 flex items-center justify-center">
        <p className="text-muted-foreground text-center">
          Start a conversation to see travel recommendations
        </p>
      </div>
    );
  }

  return (
    <div className="w-96 bg-card border-l border-border flex flex-col h-screen">
      <div className="p-6 border-b border-border">
        <h2 className="text-2xl font-bold mb-2">{travelData.destination}</h2>
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">{travelData.dates}</p>
          <p className="text-lg font-bold text-primary">
            {travelData.totalCost}
          </p>
        </div>
      </div>

      <div className="h-64 relative">
        <MapContainer
          center={travelData.coordinates}
          zoom={13}
          className="h-full w-full"
          scrollWheelZoom={false}
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <Marker position={travelData.coordinates}>
            <Popup>{travelData.destination}</Popup>
          </Marker>
        </MapContainer>
      </div>

      <Tabs defaultValue="itinerary" className="flex-1 flex flex-col">
        <TabsList className="m-4 grid w-auto grid-cols-4">
          <TabsTrigger value="itinerary">
            <MapPin className="h-4 w-4" />
          </TabsTrigger>
          <TabsTrigger value="restaurant">
            <Utensils className="h-4 w-4" />
          </TabsTrigger>
          <TabsTrigger value="hotel">
            <Hotel className="h-4 w-4" />
          </TabsTrigger>
          <TabsTrigger value="flights">
            <Plane className="h-4 w-4" />
          </TabsTrigger>
        </TabsList>

        <TabsContent value="itinerary" className="flex-1 m-0">
          <ScrollArea className="h-full px-4 pb-4">
            <div className="space-y-4">
              {travelData.itinerary.map((day) => (
                <Card key={day.day} className="shadow-card">
                  <CardHeader>
                    <CardTitle className="text-sm">Day {day.day}</CardTitle>
                    <p className="text-xs text-muted-foreground">{day.title}</p>
                  </CardHeader>
                  <CardContent>
                    <ul className="space-y-2">
                      {day.activities.map((activity, idx) => (
                        <li
                          key={idx}
                          className="text-sm flex items-start gap-2"
                        >
                          <span className="text-primary mt-1">•</span>
                          <span>{activity}</span>
                        </li>
                      ))}
                    </ul>
                  </CardContent>
                </Card>
              ))}
            </div>
          </ScrollArea>
        </TabsContent>

        <TabsContent value="restaurant" className="flex-1 px-4">
          <p className="text-sm text-muted-foreground">
            Restaurant recommendations coming soon...
          </p>
        </TabsContent>

        <TabsContent value="hotel" className="flex-1 px-4">
          <p className="text-sm text-muted-foreground">
            Hotel recommendations coming soon...
          </p>
        </TabsContent>

        <TabsContent value="flights" className="flex-1 px-4">
          <p className="text-sm text-muted-foreground">
            Flight options coming soon...
          </p>
        </TabsContent>
      </Tabs>
    </div>
  );
};
