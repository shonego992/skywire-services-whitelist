import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'uptime'
  })
  export class UptimePipe implements PipeTransform {
  
      transform(value: number): string {
        const day: number = Math.floor(value / 86400);
        value -= day * 86400;
        const hours: number = Math.floor(value/3600);
        value -= hours * 3600;
        const minutes: number = Math.floor(value / 60);
        value -= minutes * 60;
        const sec:number = Math.floor(value);
        var timeStamp:string = '';
        if (day > 0) {
            timeStamp += day + 'd ';
        }
        if (hours > 0) {
            timeStamp += hours + 'h ';
        }
        if (minutes > 0) {
            timeStamp += minutes + 'm ';
        }
        if (sec > 0) {
            timeStamp += sec + 's ';
        }
         return timeStamp;
      }
  
  }