import { Module } from '@nestjs/common';
import { SpotsService } from '@app/core/spots/spots.service';
import { LugaresController } from './lugares.controller';

@Module({
  controllers: [LugaresController],
  providers: [SpotsService],
})
export class LugaresModule { }
